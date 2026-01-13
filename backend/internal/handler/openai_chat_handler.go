// Package handler provides HTTP handlers for OpenAI Chat Completions API.
package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// OpenAIChatHandler handles OpenAI Chat Completions API requests
type OpenAIChatHandler struct {
	chatService               *service.OpenAIChatService
	geminiCompatService       *service.GeminiMessagesCompatService
	antigravityGatewayService *service.AntigravityGatewayService
	apiKeyService             *service.APIKeyService
}

// NewOpenAIChatHandler creates a new OpenAIChatHandler
func NewOpenAIChatHandler(
	chatService *service.OpenAIChatService,
	geminiCompatService *service.GeminiMessagesCompatService,
	antigravityGatewayService *service.AntigravityGatewayService,
	apiKeyService *service.APIKeyService,
) *OpenAIChatHandler {
	return &OpenAIChatHandler{
		chatService:               chatService,
		geminiCompatService:       geminiCompatService,
		antigravityGatewayService: antigravityGatewayService,
		apiKeyService:             apiKeyService,
	}
}

// ChatCompletions handles POST /v1/chat/completions
// This endpoint proxies requests to LiteLLM with optional Claude cache injection
// For Gemini models, it handles the request directly without forwarding to LiteLLM
func (h *OpenAIChatHandler) ChatCompletions(c *gin.Context) {
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.errorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
		}
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
	}

	if len(body) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
	}

	// 解析请求体获取模型信息
	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
	}

	reqModel, _ := reqBody["model"].(string)
	if reqModel == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
	}

	// 设置 ops 上下文（如果需要）
	reqStream, _ := reqBody["stream"].(bool)
	setOpsRequestContext(c, reqModel, reqStream, body)

	// 检查是否是 Gemini 模型，如果是则直接处理而不转发到 LiteLLM
	if h.isGeminiModel(reqModel) {
		h.handleGeminiRequest(c, body, reqBody)
		return
	}

	// 转发到 LiteLLM
	result, err := h.chatService.Forward(c.Request.Context(), c, body)
	if err != nil {
		log.Printf("[OpenAI Chat] Forward error: %v", err)
		// 响应可能已经被写入（流式），不再返回错误
		if !c.Writer.Written() {
			h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Failed to forward request to LiteLLM")
		}
		return
	}

	// 记录缓存使用情况
	if result != nil && (result.Usage.CacheCreationInputTokens > 0 || result.Usage.CacheReadInputTokens > 0) {
		log.Printf("[OpenAI Chat] Cache stats: model=%s creation=%d read=%d duration=%v",
			result.Model,
			result.Usage.CacheCreationInputTokens,
			result.Usage.CacheReadInputTokens,
			result.Duration)
	}
}

// isGeminiModel checks if the model is a Gemini model
func (h *OpenAIChatHandler) isGeminiModel(model string) bool {
	lower := strings.ToLower(model)
	return strings.HasPrefix(lower, "gemini-") || strings.Contains(lower, "gemini")
}

// handleGeminiRequest handles Gemini model requests directly using Antigravity service
func (h *OpenAIChatHandler) handleGeminiRequest(c *gin.Context, body []byte, reqBody map[string]any) {
	ctx := c.Request.Context()

	// 从 Authorization header 获取 API key
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Authorization header is required")
		return
	}

	// 解析 Bearer token
	apiKeyStr := strings.TrimPrefix(authHeader, "Bearer ")
	if apiKeyStr == authHeader || apiKeyStr == "" {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid authorization format")
		return
	}

	// 验证 API key
	apiKey, _, err := h.apiKeyService.ValidateKey(ctx, apiKeyStr)
	if err != nil {
		log.Printf("[OpenAI Chat Gemini] API key validation error: %v", err)
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}

	reqModel, _ := reqBody["model"].(string)

	// 选择 Antigravity 账户（用于 Gemini 模型）
	var groupID *int64
	if apiKey.GroupID != nil {
		groupID = apiKey.GroupID
	}

	// 使用 GeminiMessagesCompatService 选择账户，它支持多平台选择
	account, err := h.geminiCompatService.SelectAccountForModel(ctx, groupID, "", reqModel)
	if err != nil {
		log.Printf("[OpenAI Chat Gemini] Select account error: %v", err)
		h.errorResponse(c, http.StatusServiceUnavailable, "service_error", "No available accounts for Gemini models")
		return
	}

	// 转换 OpenAI 格式到 Claude 格式（Forward 方法内部会再转换为 Gemini 格式）
	claudeBody := h.convertOpenAIToClaude(body, reqBody)

	// 根据账户平台选择正确的服务
	var result *service.ForwardResult
	if account.Platform == service.PlatformAntigravity {
		// 使用 Antigravity 服务
		result, err = h.antigravityGatewayService.Forward(ctx, c, account, claudeBody)
	} else {
		// 使用 Gemini 服务（原生 Gemini OAuth 账户）
		result, err = h.geminiCompatService.Forward(ctx, c, account, claudeBody)
	}

	if err != nil {
		log.Printf("[OpenAI Chat Gemini] Forward error: %v", err)
		if !c.Writer.Written() {
			h.errorResponse(c, http.StatusBadGateway, "upstream_error", "Failed to forward request to Gemini")
		}
		return
	}

	if result != nil {
		log.Printf("[OpenAI Chat Gemini] Request completed: model=%s duration=%v", reqModel, result.Duration)
	}
}

// convertOpenAIToClaude converts OpenAI chat format to Claude messages format
func (h *OpenAIChatHandler) convertOpenAIToClaude(body []byte, reqBody map[string]any) []byte {
	messages, ok := reqBody["messages"].([]any)
	if !ok {
		return body
	}

	// 提取 system message
	var systemText string
	var newMessages []any

	for _, msg := range messages {
		m, ok := msg.(map[string]any)
		if !ok {
			continue
		}
		role, _ := m["role"].(string)

		if role == "system" {
			// 提取 system 内容
			content := m["content"]
			switch v := content.(type) {
			case string:
				systemText = v
			case []any:
				var parts []string
				for _, p := range v {
					if pm, ok := p.(map[string]any); ok {
						if text, ok := pm["text"].(string); ok {
							parts = append(parts, text)
						}
					}
				}
				systemText = strings.Join(parts, "\n")
			}
		} else {
			newMessages = append(newMessages, msg)
		}
	}

	// 构建 Claude 格式请求
	claudeReq := make(map[string]any)
	for k, v := range reqBody {
		claudeReq[k] = v
	}

	claudeReq["messages"] = newMessages
	if systemText != "" {
		claudeReq["system"] = systemText
	}

	result, err := json.Marshal(claudeReq)
	if err != nil {
		return body
	}
	return result
}

// errorResponse sends an OpenAI-style error response
func (h *OpenAIChatHandler) errorResponse(c *gin.Context, status int, errType string, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": message,
		},
	})
}
