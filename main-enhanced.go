package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Enhanced structures for multi-model support
type ModelRouter struct {
	ChatModelURL     string `json:"chat_model_url"`
	AnalysisModelURL string `json:"analysis_model_url"`
	CodeModelURL     string `json:"code_model_url"`
	MCPGatewayURL    string `json:"mcp_gateway_url"`
}

type TaskClassification struct {
	TaskType   string  `json:"task_type"`
	Confidence float64 `json:"confidence"`
	ModelURL   string  `json:"model_url"`
	MCPTools   []string `json:"mcp_tools"`
}

type EnhancedChatRequest struct {
	Message        string   `json:"message"`
	TaskType       string   `json:"task_type,omitempty"`
	PreferredModel string   `json:"preferred_model,omitempty"`
	EnabledTools   []string `json:"enabled_tools,omitempty"`
	SessionID      string   `json:"session_id,omitempty"`
}

type EnhancedChatResponse struct {
	Response       string             `json:"response"`
	ModelUsed      string             `json:"model_used"`
	TaskType       string             `json:"task_type"`
	ToolsUsed      []string           `json:"tools_used"`
	ProcessingTime time.Duration      `json:"processing_time"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type MCPToolClient struct {
	BaseURL string
	Tools   []string
}

// Enhanced metrics
var (
	modelSelectionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "aiwatch_model_selection_duration_seconds",
			Help: "Time spent selecting optimal model",
		},
		[]string{"task_type", "model_selected"},
	)

	mcpToolUsage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiwatch_mcp_tool_usage_total",
			Help: "Total number of MCP tool calls",
		},
		[]string{"tool_name", "status"},
	)

	multiModelRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiwatch_multi_model_requests_total",
			Help: "Total requests by model type",
		},
		[]string{"model_type", "task_type"},
	)
)

func init() {
	prometheus.MustRegister(modelSelectionDuration)
	prometheus.MustRegister(mcpToolUsage)
	prometheus.MustRegister(multiModelRequests)
}

// Enhanced AI service with multi-model support
type EnhancedAIService struct {
	ModelRouter    *ModelRouter
	MCPClient      *MCPToolClient
	Tracer         trace.Tracer
	FeatureFlags   map[string]bool
}

func NewEnhancedAIService() *EnhancedAIService {
	return &EnhancedAIService{
		ModelRouter: &ModelRouter{
			ChatModelURL:     getEnv("PRIMARY_MODEL_URL", "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
			AnalysisModelURL: getEnv("ANALYSIS_MODEL_URL", "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
			CodeModelURL:     getEnv("CODE_MODEL_URL", "http://model-runner.docker.internal/engines/llama.cpp/v1/"),
			MCPGatewayURL:    getEnv("MCP_GATEWAY_URL", "http://mcp-gateway:8811"),
		},
		MCPClient: &MCPToolClient{
			BaseURL: getEnv("MCP_GATEWAY_URL", "http://mcp-gateway:8811"),
			Tools:   []string{"web_search", "document_processor", "code_assistant"},
		},
		Tracer: otel.Tracer("aiwatch-enhanced"),
		FeatureFlags: map[string]bool{
			"multi_model_enabled": getEnv("MULTI_MODEL_ENABLED", "true") == "true",
			"mcp_tools_enabled":   getEnv("MCP_TOOLS_ENABLED", "true") == "true",
			"intelligent_routing": getEnv("INTELLIGENT_ROUTING", "true") == "true",
		},
	}
}

// Intelligent task classification
func (s *EnhancedAIService) ClassifyTask(ctx context.Context, message string) *TaskClassification {
	ctx, span := s.Tracer.Start(ctx, "classify_task")
	defer span.End()

	startTime := time.Now()
	defer func() {
		modelSelectionDuration.WithLabelValues("auto", "classifier").Observe(time.Since(startTime).Seconds())
	}()

	// Simple rule-based classification (can be enhanced with ML)
	message = strings.ToLower(message)
	
	classification := &TaskClassification{
		Confidence: 0.8,
		MCPTools:   []string{},
	}

	switch {
	case strings.Contains(message, "code") || strings.Contains(message, "function") || 
		 strings.Contains(message, "debug") || strings.Contains(message, "refactor"):
		classification.TaskType = "code"
		classification.ModelURL = s.ModelRouter.CodeModelURL
		classification.MCPTools = []string{"code_assistant", "document_processor"}
		
	case strings.Contains(message, "analyze") || strings.Contains(message, "research") || 
		 strings.Contains(message, "compare") || strings.Contains(message, "evaluate"):
		classification.TaskType = "analysis"
		classification.ModelURL = s.ModelRouter.AnalysisModelURL
		classification.MCPTools = []string{"web_research", "document_processor"}
		
	case strings.Contains(message, "search") || strings.Contains(message, "find") || 
		 strings.Contains(message, "lookup"):
		classification.TaskType = "research"
		classification.ModelURL = s.ModelRouter.AnalysisModelURL
		classification.MCPTools = []string{"web_research"}
		
	default:
		classification.TaskType = "chat"
		classification.ModelURL = s.ModelRouter.ChatModelURL
	}

	return classification
}

// Enhanced chat processing with multi-model support
func (s *EnhancedAIService) ProcessEnhancedChat(ctx context.Context, req *EnhancedChatRequest) (*EnhancedChatResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "process_enhanced_chat")
	defer span.End()

	startTime := time.Now()
	
	// Task classification
	var classification *TaskClassification
	if req.TaskType != "" {
		// Use specified task type
		classification = &TaskClassification{
			TaskType: req.TaskType,
			ModelURL: s.getModelURLByType(req.TaskType),
		}
	} else if s.FeatureFlags["intelligent_routing"] {
		// Intelligent routing
		classification = s.ClassifyTask(ctx, req.Message)
	} else {
		// Default to chat model
		classification = &TaskClassification{
			TaskType: "chat",
			ModelURL: s.ModelRouter.ChatModelURL,
		}
	}

	// Override with preferred model if specified
	if req.PreferredModel != "" {
		classification.ModelURL = s.getModelURLByType(req.PreferredModel)
		classification.TaskType = req.PreferredModel
	}

	// Track model usage
	multiModelRequests.WithLabelValues(classification.TaskType, classification.TaskType).Inc()

	// Prepare enhanced message with MCP tools if enabled
	enhancedMessage := req.Message
	var toolsUsed []string

	if s.FeatureFlags["mcp_tools_enabled"] && len(classification.MCPTools) > 0 {
		mcpResponse, err := s.UseMCPTools(ctx, req.Message, classification.MCPTools)
		if err == nil && mcpResponse != "" {
			enhancedMessage = fmt.Sprintf("%s\n\nAdditional context: %s", req.Message, mcpResponse)
			toolsUsed = classification.MCPTools
		}
	}

	// Call the selected model
	response, err := s.CallModel(ctx, classification.ModelURL, enhancedMessage)
	if err != nil {
		return nil, fmt.Errorf("model call failed: %w", err)
	}

	return &EnhancedChatResponse{
		Response:       response,
		ModelUsed:      classification.TaskType,
		TaskType:       classification.TaskType,
		ToolsUsed:      toolsUsed,
		ProcessingTime: time.Since(startTime),
		Metadata: map[string]interface{}{
			"classification_confidence": classification.Confidence,
			"model_url":                classification.ModelURL,
			"session_id":               req.SessionID,
		},
	}, nil
}

// MCP tools integration
func (s *EnhancedAIService) UseMCPTools(ctx context.Context, message string, tools []string) (string, error) {
	ctx, span := s.Tracer.Start(ctx, "use_mcp_tools")
	defer span.End()

	var results []string
	
	for _, tool := range tools {
		startTime := time.Now()
		
		result, err := s.callMCPTool(ctx, tool, message)
		
		if err != nil {
			mcpToolUsage.WithLabelValues(tool, "error").Inc()
			log.Printf("MCP tool %s failed: %v", tool, err)
			continue
		}
		
		mcpToolUsage.WithLabelValues(tool, "success").Inc()
		results = append(results, result)
		
		log.Printf("MCP tool %s completed in %v", tool, time.Since(startTime))
	}

	return strings.Join(results, "\n"), nil
}

func (s *EnhancedAIService) callMCPTool(ctx context.Context, tool string, message string) (string, error) {
	// Implementation for calling MCP tools via gateway
	// This would integrate with Docker's MCP Gateway
	url := fmt.Sprintf("%s/tools/%s", s.MCPClient.BaseURL, tool)
	
	payload := map[string]interface{}{
		"input": message,
		"tool":  tool,
	}
	
	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body), nil // Return raw response if JSON parsing fails
	}
	
	if output, ok := result["output"].(string); ok {
		return output, nil
	}
	
	return string(body), nil
}

func (s *EnhancedAIService) CallModel(ctx context.Context, modelURL string, message string) (string, error) {
	// Your existing model calling logic, enhanced
	// This is similar to your current implementation but with model URL selection
	
	payload := map[string]interface{}{
		"model": "llama3.2", // This would be dynamic based on modelURL
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
		"stream": false,
	}
	
	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, "POST", modelURL+"chat/completions", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getEnv("API_KEY", "ollama"))
	
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body), nil
	}
	
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}
	
	return string(body), nil
}

func (s *EnhancedAIService) getModelURLByType(modelType string) string {
	switch modelType {
	case "code":
		return s.ModelRouter.CodeModelURL
	case "analysis", "research":
		return s.ModelRouter.AnalysisModelURL
	default:
		return s.ModelRouter.ChatModelURL
	}
}

// HTTP Handlers
func (s *EnhancedAIService) handleEnhancedChat(c *gin.Context) {
	var req EnhancedChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	response, err := s.ProcessEnhancedChat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

func (s *EnhancedAIService) handleModelCapabilities(c *gin.Context) {
	capabilities := map[string]interface{}{
		"available_models": []string{"chat", "analysis", "code"},
		"mcp_tools":       s.MCPClient.Tools,
		"feature_flags":   s.FeatureFlags,
		"model_urls": map[string]string{
			"chat":     s.ModelRouter.ChatModelURL,
			"analysis": s.ModelRouter.AnalysisModelURL,
			"code":     s.ModelRouter.CodeModelURL,
		},
	}
	
	c.JSON(http.StatusOK, capabilities)
}

// Utility functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize enhanced service
	service := NewEnhancedAIService()
	
	// Setup Gin router
	r := gin.Default()
	
	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3002"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api/v1")
	{
		// Enhanced chat endpoint
		api.POST("/chat/enhanced", service.handleEnhancedChat)
		
		// Model capabilities endpoint
		api.GET("/capabilities", service.handleModelCapabilities)
		
		// Backward compatibility - your existing chat endpoint
		api.POST("/chat", func(c *gin.Context) {
			// Convert to enhanced request for backward compatibility
			var oldReq map[string]interface{}
			if err := c.ShouldBindJSON(&oldReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			
			enhancedReq := &EnhancedChatRequest{
				Message: oldReq["message"].(string),
			}
			
			response, err := service.ProcessEnhancedChat(c.Request.Context(), enhancedReq)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			// Return in old format for compatibility
			c.JSON(http.StatusOK, gin.H{
				"response": response.Response,
			})
		})
	}

	// Health and metrics endpoints (your existing ones)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "enhanced": true})
	})
	
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Enhanced AIWatch backend starting on port %s", port)
	log.Printf("Multi-model enabled: %v", service.FeatureFlags["multi_model_enabled"])
	log.Printf("MCP tools enabled: %v", service.FeatureFlags["mcp_tools_enabled"])
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
