package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RedisTimeSeriesService provides time-series analytics using Redis TimeSeries
type RedisTimeSeriesService struct {
	redis *redis.Client
	ctx   context.Context
	
	// Prometheus metrics
	timeSeriesOperations *prometheus.CounterVec
	timeSeriesLatency    *prometheus.HistogramVec
}

// TimeSeriesMetric represents a time-series data point
type TimeSeriesMetric struct {
	Key       string                 `json:"key"`
	Timestamp int64                  `json:"timestamp"`
	Value     float64                `json:"value"`
	Labels    map[string]interface{} `json:"labels"`
}

// TimeSeriesQuery represents a query for time-series data
type TimeSeriesQuery struct {
	Key       string `json:"key"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Aggregation string `json:"aggregation,omitempty"` // avg, sum, min, max, count
	BucketDuration int64 `json:"bucket_duration,omitempty"` // in milliseconds
}

// TimeSeriesResponse represents the response for time-series queries
type TimeSeriesResponse struct {
	Key    string      `json:"key"`
	Data   []DataPoint `json:"data"`
	Labels map[string]interface{} `json:"labels"`
}

type DataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// NewRedisTimeSeriesService creates a new time-series service
func NewRedisTimeSeriesService(redisAddr, redisPassword string, redisDB int) *RedisTimeSeriesService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize Prometheus metrics
	timeSeriesOperations := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_timeseries_operations_total",
			Help: "Total number of time-series operations",
		},
		[]string{"operation", "status"},
	)

	timeSeriesLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_timeseries_operation_duration_seconds",
			Help:    "Time-series operation latency",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"operation"},
	)

	// Register metrics
	prometheus.MustRegister(timeSeriesOperations, timeSeriesLatency)

	service := &RedisTimeSeriesService{
		redis:                rdb,
		ctx:                  ctx,
		timeSeriesOperations: timeSeriesOperations,
		timeSeriesLatency:    timeSeriesLatency,
	}

	// Initialize time-series keys
	service.initializeTimeSeries()

	return service
}

// initializeTimeSeries creates time-series keys with appropriate retention and labels
func (ts *RedisTimeSeriesService) initializeTimeSeries() {
	timeSeries := map[string]map[string]interface{}{
		"metrics:tokens:input_rate": {
			"RETENTION": 86400000, // 24 hours in milliseconds
			"LABELS": map[string]string{
				"metric_type": "token_rate",
				"direction":   "input",
			},
		},
		"metrics:tokens:output_rate": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "token_rate",
				"direction":   "output",
			},
		},
		"metrics:users:active_5m": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "user_activity",
				"window":      "5m",
			},
		},
		"metrics:users:active_1h": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "user_activity",
				"window":      "1h",
			},
		},
		"metrics:response_time:p95": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "response_time",
				"percentile":  "95",
			},
		},
		"metrics:response_time:p99": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "response_time",
				"percentile":  "99",
			},
		},
		"metrics:error_rate": {
			"RETENTION": 86400000,
			"LABELS": map[string]string{
				"metric_type": "error_rate",
			},
		},
		"metrics:memory:redis_used": {
			"RETENTION": 604800000, // 7 days
			"LABELS": map[string]string{
				"metric_type": "memory",
				"component":   "redis",
			},
		},
		"metrics:cpu:usage": {
			"RETENTION": 604800000,
			"LABELS": map[string]string{
				"metric_type": "system",
				"component":   "cpu",
			},
		},
	}

	for key, config := range timeSeries {
		// Create time-series with labels and retention
		args := []interface{}{"TS.CREATE", key}
		
		if retention, ok := config["RETENTION"]; ok {
			args = append(args, "RETENTION", retention)
		}
		
		if labels, ok := config["LABELS"].(map[string]string); ok {
			args = append(args, "LABELS")
			for labelKey, labelValue := range labels {
				args = append(args, labelKey, labelValue)
			}
		}

		// Execute create command (ignore if already exists)
		err := ts.redis.Do(ts.ctx, args...).Err()
		if err != nil && err.Error() != "TSDB: key already exists" {
			log.Printf("Warning: Failed to create time-series %s: %v", key, err)
		}
	}

	log.Println("Time-series initialization completed")
}

// AddDataPoint adds a data point to a time-series
func (ts *RedisTimeSeriesService) AddDataPoint(key string, timestamp int64, value float64) error {
	start := time.Now()
	defer func() {
		ts.timeSeriesLatency.WithLabelValues("add").Observe(time.Since(start).Seconds())
	}()

	// If timestamp is 0, use current time
	if timestamp == 0 {
		timestamp = time.Now().UnixMilli()
	}

	err := ts.redis.Do(ts.ctx, "TS.ADD", key, timestamp, value).Err()
	
	status := "success"
	if err != nil {
		status = "error"
	}
	ts.timeSeriesOperations.WithLabelValues("add", status).Inc()

	return err
}

// QueryRange queries time-series data for a range
func (ts *RedisTimeSeriesService) QueryRange(query TimeSeriesQuery) (*TimeSeriesResponse, error) {
	start := time.Now()
	defer func() {
		ts.timeSeriesLatency.WithLabelValues("query_range").Observe(time.Since(start).Seconds())
	}()

	args := []interface{}{"TS.RANGE", query.Key, query.StartTime, query.EndTime}

	// Add aggregation if specified
	if query.Aggregation != "" && query.BucketDuration > 0 {
		args = append(args, "AGGREGATION", query.Aggregation, query.BucketDuration)
	}

	result, err := ts.redis.Do(ts.ctx, args...).Result()
	
	status := "success"
	if err != nil {
		status = "error"
	}
	ts.timeSeriesOperations.WithLabelValues("query_range", status).Inc()

	if err != nil {
		return nil, err
	}

	// Parse result
	response := &TimeSeriesResponse{
		Key:  query.Key,
		Data: []DataPoint{},
	}

	// Parse Redis TimeSeries response format
	if resultSlice, ok := result.([]interface{}); ok {
		for _, item := range resultSlice {
			if itemSlice, ok := item.([]interface{}); ok && len(itemSlice) == 2 {
				if timestamp, ok := itemSlice[0].(int64); ok {
					if valueStr, ok := itemSlice[1].(string); ok {
						if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
							response.Data = append(response.Data, DataPoint{
								Timestamp: timestamp,
								Value:     value,
							})
						}
					}
				}
			}
		}
	}

	return response, nil
}

// QueryMultiRange queries multiple time-series
func (ts *RedisTimeSeriesService) QueryMultiRange(queries []TimeSeriesQuery) (map[string]*TimeSeriesResponse, error) {
	start := time.Now()
	defer func() {
		ts.timeSeriesLatency.WithLabelValues("query_multi_range").Observe(time.Since(start).Seconds())
	}()

	results := make(map[string]*TimeSeriesResponse)
	
	for _, query := range queries {
		response, err := ts.QueryRange(query)
		if err != nil {
			ts.timeSeriesOperations.WithLabelValues("query_multi_range", "error").Inc()
			return nil, fmt.Errorf("failed to query %s: %v", query.Key, err)
		}
		results[query.Key] = response
	}

	ts.timeSeriesOperations.WithLabelValues("query_multi_range", "success").Inc()
	return results, nil
}

// GetLatestValue gets the latest value for a time-series
func (ts *RedisTimeSeriesService) GetLatestValue(key string) (*DataPoint, error) {
	start := time.Now()
	defer func() {
		ts.timeSeriesLatency.WithLabelValues("get_latest").Observe(time.Since(start).Seconds())
	}()

	result, err := ts.redis.Do(ts.ctx, "TS.GET", key).Result()
	
	status := "success"
	if err != nil {
		status = "error"
	}
	ts.timeSeriesOperations.WithLabelValues("get_latest", status).Inc()

	if err != nil {
		return nil, err
	}

	// Parse result
	if resultSlice, ok := result.([]interface{}); ok && len(resultSlice) == 2 {
		if timestamp, ok := resultSlice[0].(int64); ok {
			if valueStr, ok := resultSlice[1].(string); ok {
				if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
					return &DataPoint{
						Timestamp: timestamp,
						Value:     value,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("invalid response format")
}

// UpdateMetricsFromRedis updates time-series from current Redis analytics data
func (ts *RedisTimeSeriesService) UpdateMetricsFromRedis() error {
	timestamp := time.Now().UnixMilli()

	// Get active users
	activeUsers5m, _ := ts.redis.SCard(ts.ctx, "users:active:5m").Result()
	activeUsers1h, _ := ts.redis.SCard(ts.ctx, "users:active:1h").Result()

	// Add to time-series
	ts.AddDataPoint("metrics:users:active_5m", timestamp, float64(activeUsers5m))
	ts.AddDataPoint("metrics:users:active_1h", timestamp, float64(activeUsers1h))

	// Get token rates (approximate from recent data)
	inputTokens, _ := ts.redis.Get(ts.ctx, "tokens:input:count").Float64()
	outputTokens, _ := ts.redis.Get(ts.ctx, "tokens:output:count").Float64()

	ts.AddDataPoint("metrics:tokens:input_rate", timestamp, inputTokens)
	ts.AddDataPoint("metrics:tokens:output_rate", timestamp, outputTokens)

	// Get error rate
	errorCount, _ := ts.redis.Get(ts.ctx, "errors:total:count").Float64()
	ts.AddDataPoint("metrics:error_rate", timestamp, errorCount)

	return nil
}

// StartMetricsCollection starts background metrics collection
func (ts *RedisTimeSeriesService) StartMetricsCollection() {
	ticker := time.NewTicker(30 * time.Second) // Collect every 30 seconds
	go func() {
		for range ticker.C {
			if err := ts.UpdateMetricsFromRedis(); err != nil {
				log.Printf("Error updating time-series metrics: %v", err)
			}
		}
	}()
}

// HTTP Handlers

func (ts *RedisTimeSeriesService) queryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var query TimeSeriesQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	response, err := ts.QueryRange(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (ts *RedisTimeSeriesService) multiQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var queries []TimeSeriesQuery
	if err := json.NewDecoder(r.Body).Decode(&queries); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	responses, err := ts.QueryMultiRange(queries)
	if err != nil {
		http.Error(w, fmt.Sprintf("Multi-query failed: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(responses)
}

func (ts *RedisTimeSeriesService) latestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	dataPoint, err := ts.GetLatestValue(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get latest value: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dataPoint)
}

func (ts *RedisTimeSeriesService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "redis-timeseries"})
}

func main() {
	// Get configuration from environment
	redisAddr := getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))
	port := getEnvOrDefault("TIMESERIES_PORT", "8082")

	log.Printf("Starting Redis TimeSeries Service on port %s", port)
	log.Printf("Connecting to Redis at %s", redisAddr)

	// Create time-series service
	service := NewRedisTimeSeriesService(redisAddr, redisPassword, redisDB)

	// Start background metrics collection
	service.StartMetricsCollection()

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/query", service.queryHandler)
	mux.HandleFunc("/multi-query", service.multiQueryHandler)
	mux.HandleFunc("/latest", service.latestHandler)
	mux.HandleFunc("/health", service.healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Redis TimeSeries Service running on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
