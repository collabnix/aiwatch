package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// TokenCaptureService handles Redis-based token tracking
type TokenCaptureService struct {
	client *redis.Client
	ctx    context.Context
}

// TokenMetrics represents the detailed token usage data
type TokenMetrics struct {
	RequestID           string  `json:"request_id"`
	SessionID           string  `json:"session_id"`
	UserID              string  `json:"user_id"`
	Timestamp           int64   `json:"timestamp"`
	InputTokens         int     `json:"input_tokens"`
	OutputTokens        int     `json:"output_tokens"`
	ResponseTimeMs      float64 `json:"response_time_ms"`
	FirstTokenLatencyMs float64 `json:"first_token_latency_ms"`
	ModelUsed           string  `json:"model_used"`
	PromptLength        int     `json:"prompt_length"`
	ResponseLength      int     `json:"response_length"`
	Status              string  `json:"status"`
}

// NewTokenCaptureService creates a new token capture service
func NewTokenCaptureService(redisAddr, redisPassword string, redisDB int) *TokenCaptureService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")

	return &TokenCaptureService{
		client: rdb,
		ctx:    ctx,
	}
}

// GenerateRequestID creates a unique request ID
func (tcs *TokenCaptureService) GenerateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "req_" + hex.EncodeToString(bytes)
}

// GenerateSessionID creates a unique session ID
func (tcs *TokenCaptureService) GenerateSessionID() string {
	return "sess_" + uuid.New().String()[:8]
}

// ExtractUserID extracts user ID from request
func (tcs *TokenCaptureService) ExtractUserID(r *http.Request) string {
	// Check for user ID in headers
	if userID := r.Header.Get("X-User-ID"); userID != "" {
		return userID
	}
	
	// Check for user ID in query params
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		return userID
	}
	
	// Check for session cookie
	if cookie, err := r.Cookie("user_session"); err == nil {
		return cookie.Value
	}
	
	// Default to IP-based user ID for demo
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return "user_" + strings.Replace(ip, ":", "_", -1)
}

// GetOrCreateSession retrieves existing session or creates new one
func (tcs *TokenCaptureService) GetOrCreateSession(userID string, modelUsed string) (string, error) {
	// Check for active session for this user
	activeSessionsKey := "sessions:active"
	activeSessions, err := tcs.client.SMembers(tcs.ctx, activeSessionsKey).Result()
	if err != nil {
		return "", err
	}

	// Look for existing active session for this user
	for _, sessionID := range activeSessions {
		sessionKey := fmt.Sprintf("session:%s:tokens", sessionID)
		sessionUserID, err := tcs.client.HGet(tcs.ctx, sessionKey, "user_id").Result()
		if err != nil {
			continue
		}
		
		if sessionUserID == userID {
			// Check if session is still active (within 30 minutes)
			lastActivityStr, err := tcs.client.HGet(tcs.ctx, sessionKey, "last_activity").Result()
			if err != nil {
				continue
			}
			
			lastActivity, err := time.Parse(time.RFC3339, lastActivityStr)
			if err != nil {
				continue
			}
			
			if time.Since(lastActivity) < 30*time.Minute {
				// Update last activity
				tcs.client.HSet(tcs.ctx, sessionKey, "last_activity", time.Now().Format(time.RFC3339))
				return sessionID, nil
			}
		}
	}

	// Create new session
	sessionID := tcs.GenerateSessionID()
	sessionKey := fmt.Sprintf("session:%s:tokens", sessionID)
	
	now := time.Now()
	sessionData := map[string]interface{}{
		"user_id":              userID,
		"start_time":           now.Format(time.RFC3339),
		"last_activity":        now.Format(time.RFC3339),
		"total_input_tokens":   0,
		"total_output_tokens":  0,
		"request_count":        0,
		"model_used":           modelUsed,
		"avg_response_time":    0.0,
		"status":               "active",
	}

	err = tcs.client.HMSet(tcs.ctx, sessionKey, sessionData).Err()
	if err != nil {
		return "", err
	}

	// Add to active sessions
	tcs.client.SAdd(tcs.ctx, activeSessionsKey, sessionID)
	
	// Set TTL for session (30 days)
	tcs.client.Expire(tcs.ctx, sessionKey, 30*24*time.Hour)

	return sessionID, nil
}

// CaptureTokenMetrics stores comprehensive token metrics in Redis
func (tcs *TokenCaptureService) CaptureTokenMetrics(metrics TokenMetrics) error {
	now := time.Now()
	metrics.Timestamp = now.Unix()

	// 1. Store request-level metrics
	requestKey := fmt.Sprintf("request:%s:tokens", metrics.RequestID)
	requestData := map[string]interface{}{
		"session_id":             metrics.SessionID,
		"user_id":                metrics.UserID,
		"timestamp":              metrics.Timestamp,
		"input_tokens":           metrics.InputTokens,
		"output_tokens":          metrics.OutputTokens,
		"response_time_ms":       metrics.ResponseTimeMs,
		"first_token_latency_ms": metrics.FirstTokenLatencyMs,
		"model_used":             metrics.ModelUsed,
		"prompt_length":          metrics.PromptLength,
		"response_length":        metrics.ResponseLength,
		"status":                 metrics.Status,
	}

	err := tcs.client.HMSet(tcs.ctx, requestKey, requestData).Err()
	if err != nil {
		return fmt.Errorf("failed to store request metrics: %v", err)
	}
	
	// Set TTL for request data (7 days)
	tcs.client.Expire(tcs.ctx, requestKey, 7*24*time.Hour)

	// 2. Update session metrics
	err = tcs.updateSessionMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to update session metrics: %v", err)
	}

	// 3. Update user metrics
	err = tcs.updateUserMetrics(metrics)
	if err != nil {
		return fmt.Errorf("failed to update user metrics: %v", err)
	}

	// 4. Update time-series data
	err = tcs.updateTimeSeriesData(metrics, now)
	if err != nil {
		return fmt.Errorf("failed to update time-series data: %v", err)
	}

	// 5. Update model usage statistics
	err = tcs.updateModelUsage(metrics)
	if err != nil {
		return fmt.Errorf("failed to update model usage: %v", err)
	}

	// 6. Update real-time activity tracking
	err = tcs.updateRealTimeActivity(metrics.UserID, metrics.SessionID)
	if err != nil {
		return fmt.Errorf("failed to update real-time activity: %v", err)
	}

	return nil
}

// Helper methods for updating different metric types
func (tcs *TokenCaptureService) updateSessionMetrics(metrics TokenMetrics) error {
	sessionKey := fmt.Sprintf("session:%s:tokens", metrics.SessionID)
	
	// Get current session data
	currentInputTokens, _ := tcs.client.HGet(tcs.ctx, sessionKey, "total_input_tokens").Int()
	currentOutputTokens, _ := tcs.client.HGet(tcs.ctx, sessionKey, "total_output_tokens").Int()
	currentRequestCount, _ := tcs.client.HGet(tcs.ctx, sessionKey, "request_count").Int()
	currentAvgResponseTime, _ := tcs.client.HGet(tcs.ctx, sessionKey, "avg_response_time").Float64()

	// Calculate new averages
	newRequestCount := currentRequestCount + 1
	newAvgResponseTime := ((currentAvgResponseTime * float64(currentRequestCount)) + metrics.ResponseTimeMs) / float64(newRequestCount)

	// Update session data
	sessionUpdates := map[string]interface{}{
		"total_input_tokens":  currentInputTokens + metrics.InputTokens,
		"total_output_tokens": currentOutputTokens + metrics.OutputTokens,
		"request_count":       newRequestCount,
		"avg_response_time":   newAvgResponseTime,
		"last_activity":       time.Unix(metrics.Timestamp, 0).Format(time.RFC3339),
	}

	return tcs.client.HMSet(tcs.ctx, sessionKey, sessionUpdates).Err()
}

func (tcs *TokenCaptureService) updateUserMetrics(metrics TokenMetrics) error {
	userKey := fmt.Sprintf("user:%s:tokens", metrics.UserID)
	
	// Get current user data
	currentInputTokens, _ := tcs.client.HGet(tcs.ctx, userKey, "total_input_tokens").Int()
	currentOutputTokens, _ := tcs.client.HGet(tcs.ctx, userKey, "total_output_tokens").Int()
	currentRequests, _ := tcs.client.HGet(tcs.ctx, userKey, "total_requests").Int()

	// Calculate new values
	newTotalInputTokens := currentInputTokens + metrics.InputTokens
	newTotalOutputTokens := currentOutputTokens + metrics.OutputTokens
	newTotalRequests := currentRequests + 1
	newAvgTokensPerRequest := float64(newTotalInputTokens+newTotalOutputTokens) / float64(newTotalRequests)

	// Check if this is the first time we see this user
	firstSeen, err := tcs.client.HGet(tcs.ctx, userKey, "first_seen").Result()
	if err == redis.Nil {
		firstSeen = time.Unix(metrics.Timestamp, 0).Format(time.RFC3339)
	}

	userUpdates := map[string]interface{}{
		"total_input_tokens":      newTotalInputTokens,
		"total_output_tokens":     newTotalOutputTokens,
		"total_requests":          newTotalRequests,
		"avg_tokens_per_request":  newAvgTokensPerRequest,
		"first_seen":              firstSeen,
		"last_seen":               time.Unix(metrics.Timestamp, 0).Format(time.RFC3339),
	}

	return tcs.client.HMSet(tcs.ctx, userKey, userUpdates).Err()
}

func (tcs *TokenCaptureService) updateTimeSeriesData(metrics TokenMetrics, timestamp time.Time) error {
	// Hourly data for user
	hourlyKey := fmt.Sprintf("user:%s:tokens:hourly:%s", 
		metrics.UserID, 
		timestamp.Format("2006-01-02-15"))
	
	minute := timestamp.Minute()
	memberData := fmt.Sprintf("input:%d:output:%d", metrics.InputTokens, metrics.OutputTokens)
	
	err := tcs.client.ZAdd(tcs.ctx, hourlyKey, &redis.Z{
		Score:  float64(minute),
		Member: memberData,
	}).Err()
	if err != nil {
		return err
	}
	
	// Set TTL for hourly data (90 days)
	tcs.client.Expire(tcs.ctx, hourlyKey, 90*24*time.Hour)

	return nil
}

func (tcs *TokenCaptureService) updateModelUsage(metrics TokenMetrics) error {
	modelKey := fmt.Sprintf("model:%s:usage", metrics.ModelUsed)
	
	// Get current model data
	currentRequests, _ := tcs.client.HGet(tcs.ctx, modelKey, "total_requests").Int()
	currentInputTokens, _ := tcs.client.HGet(tcs.ctx, modelKey, "total_input_tokens").Int()
	currentOutputTokens, _ := tcs.client.HGet(tcs.ctx, modelKey, "total_output_tokens").Int()
	currentAvgResponseTime, _ := tcs.client.HGet(tcs.ctx, modelKey, "avg_response_time").Float64()

	// Calculate new values
	newRequests := currentRequests + 1
	newAvgResponseTime := ((currentAvgResponseTime * float64(currentRequests)) + metrics.ResponseTimeMs) / float64(newRequests)

	modelUpdates := map[string]interface{}{
		"total_requests":      newRequests,
		"total_input_tokens":  currentInputTokens + metrics.InputTokens,
		"total_output_tokens": currentOutputTokens + metrics.OutputTokens,
		"avg_response_time":   newAvgResponseTime,
		"last_used":           time.Unix(metrics.Timestamp, 0).Format(time.RFC3339),
	}

	return tcs.client.HMSet(tcs.ctx, modelKey, modelUpdates).Err()
}

func (tcs *TokenCaptureService) updateRealTimeActivity(userID, sessionID string) error {
	// Add to active sessions
	tcs.client.SAdd(tcs.ctx, "sessions:active", sessionID)
	
	// Add to active users for different time windows
	timeWindows := map[string]time.Duration{
		"5m":  5 * time.Minute,
		"15m": 15 * time.Minute,
		"1h":  1 * time.Hour,
		"24h": 24 * time.Hour,
	}

	for window, duration := range timeWindows {
		key := fmt.Sprintf("users:active:%s", window)
		tcs.client.SAdd(tcs.ctx, key, userID)
		tcs.client.Expire(tcs.ctx, key, duration)
	}

	return nil
}
