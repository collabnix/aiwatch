package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TokenAnalyticsService provides real-time analytics from Redis data
type TokenAnalyticsService struct {
	redis  *redis.Client
	ctx    context.Context
	
	// Prometheus metrics
	activeUsersGauge     *prometheus.GaugeVec
	activeSessionsGauge  prometheus.Gauge
	tokenRateGauge       *prometheus.GaugeVec
	userTokensCounter    *prometheus.CounterVec
	modelUsageGauge      *prometheus.GaugeVec
	responseTimeHist     *prometheus.HistogramVec
	errorRateGauge       *prometheus.GaugeVec
}

// AnalyticsResponse represents the API response for analytics data
type AnalyticsResponse struct {
	ActiveUsers5m     int64                  `json:"active_users_5m"`
	ActiveUsers1h     int64                  `json:"active_users_1h"`
	ActiveSessions    int64                  `json:"active_sessions"`
	TokenRates        map[string]float64     `json:"token_rates"`
	TopUsers          []UserStats            `json:"top_users"`
	ModelUsage        map[string]ModelStats  `json:"model_usage"`
	ResponseTimeP95   float64                `json:"response_time_p95"`
	ResponseTimeP99   float64                `json:"response_time_p99"`
	ErrorRate         float64                `json:"error_rate"`
	Timestamp         int64                  `json:"timestamp"`
}

type UserStats struct {
	UserID              string  `json:"user_id"`
	TotalInputTokens    int64   `json:"total_input_tokens"`
	TotalOutputTokens   int64   `json:"total_output_tokens"`
	TotalSessions       int64   `json:"total_sessions"`
	AvgTokensPerRequest float64 `json:"avg_tokens_per_request"`
	LastSeen            string  `json:"last_seen"`
}

type ModelStats struct {
	TotalRequests      int64   `json:"total_requests"`
	TotalInputTokens   int64   `json:"total_input_tokens"`
	TotalOutputTokens  int64   `json:"total_output_tokens"`
	AvgResponseTime    float64 `json:"avg_response_time"`
	AvgTokensPerSecond float64 `json:"avg_tokens_per_second"`
}

// NewTokenAnalyticsService creates a new analytics service
func NewTokenAnalyticsService(redisAddr, redisPassword string, redisDB int) *TokenAnalyticsService {
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
	activeUsersGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_analytics_active_users",
			Help: "Number of active users in different time windows",
		},
		[]string{"window"},
	)

	activeSessionsGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "token_analytics_active_sessions",
			Help: "Number of currently active sessions",
		},
	)

	tokenRateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_analytics_token_rate",
			Help: "Token processing rate by direction",
		},
		[]string{"direction", "window"},
	)

	userTokensCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_analytics_user_tokens_total",
			Help: "Total tokens processed per user",
		},
		[]string{"user_id", "direction"},
	)

	modelUsageGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_analytics_model_usage",
			Help: "Model usage statistics",
		},
		[]string{"model", "metric"},
	)

	responseTimeHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "token_analytics_response_time_seconds",
			Help: "Response time distribution",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 20, 30, 60},
		},
		[]string{"model"},
	)

	errorRateGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "token_analytics_error_rate",
			Help: "Error rate by type",
		},
		[]string{"error_type"},
	)

	// Register metrics
	prometheus.MustRegister(
		activeUsersGauge,
		activeSessionsGauge,
		tokenRateGauge,
		userTokensCounter,
		modelUsageGauge,
		responseTimeHist,
		errorRateGauge,
	)

	service := &TokenAnalyticsService{
		redis:               rdb,
		ctx:                 ctx,
		activeUsersGauge:    activeUsersGauge,
		activeSessionsGauge: activeSessionsGauge,
		tokenRateGauge:      tokenRateGauge,
		userTokensCounter:   userTokensCounter,
		modelUsageGauge:     modelUsageGauge,
		responseTimeHist:    responseTimeHist,
		errorRateGauge:      errorRateGauge,
	}

	// Start background metrics collection
	go service.collectMetricsPeriodically()

	return service
}

// collectMetricsPeriodically updates Prometheus metrics from Redis data
func (tas *TokenAnalyticsService) collectMetricsPeriodically() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		tas.updatePrometheusMetrics()
	}
}

// updatePrometheusMetrics reads from Redis and updates Prometheus metrics
func (tas *TokenAnalyticsService) updatePrometheusMetrics() {
	// Update active users
	windows := []string{"5m", "15m", "1h", "24h"}
	for _, window := range windows {
		key := fmt.Sprintf("users:active:%s", window)
		count, err := tas.redis.SCard(tas.ctx, key).Result()
		if err == nil {
			tas.activeUsersGauge.WithLabelValues(window).Set(float64(count))
		}
	}

	// Update active sessions
	activeSessions, err := tas.redis.SCard(tas.ctx, "sessions:active").Result()
	if err == nil {
		tas.activeSessionsGauge.Set(float64(activeSessions))
	}

	// Update model usage statistics
	models, err := tas.redis.Keys(tas.ctx, "model:*:usage").Result()
	if err == nil {
		for _, modelKey := range models {
			modelName := strings.Split(modelKey, ":")[1]
			
			totalRequests, _ := tas.redis.HGet(tas.ctx, modelKey, "total_requests").Float64()
			totalInputTokens, _ := tas.redis.HGet(tas.ctx, modelKey, "total_input_tokens").Float64()
			totalOutputTokens, _ := tas.redis.HGet(tas.ctx, modelKey, "total_output_tokens").Float64()
			avgResponseTime, _ := tas.redis.HGet(tas.ctx, modelKey, "avg_response_time").Float64()

			tas.modelUsageGauge.WithLabelValues(modelName, "requests").Set(totalRequests)
			tas.modelUsageGauge.WithLabelValues(modelName, "input_tokens").Set(totalInputTokens)
			tas.modelUsageGauge.WithLabelValues(modelName, "output_tokens").Set(totalOutputTokens)
			tas.modelUsageGauge.WithLabelValues(modelName, "avg_response_time").Set(avgResponseTime)
		}
	}

	// Update error rates
	errorTypes := []string{"timeout", "error", "rate_limit"}
	for _, errorType := range errorTypes {
		key := fmt.Sprintf("errors:%s:count", errorType)
		count, err := tas.redis.Get(tas.ctx, key).Float64()
		if err == nil {
			tas.errorRateGauge.WithLabelValues(errorType).Set(count)
		}
	}
}

// GetAnalytics returns comprehensive analytics data
func (tas *TokenAnalyticsService) GetAnalytics() (*AnalyticsResponse, error) {
	response := &AnalyticsResponse{
		Timestamp: time.Now().Unix(),
	}

	// Get active users and sessions
	response.ActiveUsers5m, _ = tas.redis.SCard(tas.ctx, "users:active:5m").Result()
	response.ActiveUsers1h, _ = tas.redis.SCard(tas.ctx, "users:active:1h").Result()
	response.ActiveSessions, _ = tas.redis.SCard(tas.ctx, "sessions:active").Result()

	// Get token rates
	response.TokenRates = make(map[string]float64)
	response.TokenRates["input_per_minute"] = 0.0
	response.TokenRates["output_per_minute"] = 0.0

	// Get top users
	topUsers, err := tas.getTopUsers(10)
	if err == nil {
		response.TopUsers = topUsers
	}

	// Get model usage
	modelUsage, err := tas.getModelUsage()
	if err == nil {
		response.ModelUsage = modelUsage
	}

	return response, nil
}

// getTopUsers retrieves top users by token usage
func (tas *TokenAnalyticsService) getTopUsers(limit int) ([]UserStats, error) {
	userKeys, err := tas.redis.Keys(tas.ctx, "user:*:tokens").Result()
	if err != nil {
		return nil, err
	}

	var users []UserStats
	for _, key := range userKeys {
		userID := strings.Split(key, ":")[1]
		
		userData, err := tas.redis.HGetAll(tas.ctx, key).Result()
		if err != nil {
			continue
		}

		inputTokens, _ := strconv.ParseInt(userData["total_input_tokens"], 10, 64)
		outputTokens, _ := strconv.ParseInt(userData["total_output_tokens"], 10, 64)
		totalRequests, _ := strconv.ParseInt(userData["total_requests"], 10, 64)
		avgTokensPerRequest, _ := strconv.ParseFloat(userData["avg_tokens_per_request"], 64)

		users = append(users, UserStats{
			UserID:              userID,
			TotalInputTokens:    inputTokens,
			TotalOutputTokens:   outputTokens,
			TotalSessions:       totalRequests, // Approximation
			AvgTokensPerRequest: avgTokensPerRequest,
			LastSeen:            userData["last_seen"],
		})
	}

	// Limit results
	if len(users) > limit {
		users = users[:limit]
	}

	return users, nil
}

// getModelUsage retrieves model usage statistics
func (tas *TokenAnalyticsService) getModelUsage() (map[string]ModelStats, error) {
	modelKeys, err := tas.redis.Keys(tas.ctx, "model:*:usage").Result()
	if err != nil {
		return nil, err
	}

	usage := make(map[string]ModelStats)
	for _, key := range modelKeys {
		modelName := strings.Split(key, ":")[1]
		
		modelData, err := tas.redis.HGetAll(tas.ctx, key).Result()
		if err != nil {
			continue
		}

		totalRequests, _ := strconv.ParseInt(modelData["total_requests"], 10, 64)
		totalInputTokens, _ := strconv.ParseInt(modelData["total_input_tokens"], 10, 64)
		totalOutputTokens, _ := strconv.ParseInt(modelData["total_output_tokens"], 10, 64)
		avgResponseTime, _ := strconv.ParseFloat(modelData["avg_response_time"], 64)

		usage[modelName] = ModelStats{
			TotalRequests:      totalRequests,
			TotalInputTokens:   totalInputTokens,
			TotalOutputTokens:  totalOutputTokens,
			AvgResponseTime:    avgResponseTime,
			AvgTokensPerSecond: 0.0, // Calculate if needed
		}
	}

	return usage, nil
}

// HTTP handlers
func (tas *TokenAnalyticsService) analyticsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	analytics, err := tas.GetAnalytics()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get analytics: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(analytics)
}

func (tas *TokenAnalyticsService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// Get configuration from environment
	redisAddr := getEnvOrDefault("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnvOrDefault("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))
	port := getEnvOrDefault("ANALYTICS_PORT", "8081")

	log.Printf("Starting Token Analytics Service on port %s", port)
	log.Printf("Connecting to Redis at %s", redisAddr)

	// Create analytics service
	service := NewTokenAnalyticsService(redisAddr, redisPassword, redisDB)

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/analytics", service.analyticsHandler)
	mux.HandleFunc("/health", service.healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	// Start server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Token Analytics Service running on :%s", port)
	log.Fatal(server.ListenAndServe())
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
