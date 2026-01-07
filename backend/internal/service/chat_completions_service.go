package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"

	"github.com/gin-gonic/gin"
)

// ChatCompletionsService 处理 OpenAI Chat Completions API
type ChatCompletionsService struct {
	gatewayService      *GatewayService
	accountRepo         AccountRepository
	groupRepo           GroupRepository
	cache               GatewayCache
	cfg                 *config.Config
	billingService      *BillingService
	billingCacheService *BillingCacheService
	concurrencyService  *ConcurrencyService
	httpUpstream        HTTPUpstream
	identityService     *IdentityService
}

// NewChatCompletionsService 创建 ChatCompletionsService
func NewChatCompletionsService(
	gatewayService *GatewayService,
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	cache GatewayCache,
	cfg *config.Config,
	billingService *BillingService,
	billingCacheService *BillingCacheService,
	concurrencyService *ConcurrencyService,
	httpUpstream HTTPUpstream,
	identityService *IdentityService,
) *ChatCompletionsService {
	return &ChatCompletionsService{
		gatewayService:      gatewayService,
		accountRepo:         accountRepo,
		groupRepo:           groupRepo,
		cache:               cache,
		cfg:                 cfg,
		billingService:      billingService,
		billingCacheService: billingCacheService,
		concurrencyService:  concurrencyService,
		httpUpstream:        httpUpstream,
		identityService:     identityService,
	}
}

// ChatForwardResult Chat 请求转发结果
type ChatForwardResult struct {
	RequestID    string
	Model        string
	Stream       bool
	Duration     time.Duration
	FirstTokenMs *int
	Usage        *openai.ChatUsage
}

// GenerateSessionHash 从 Chat 请求生成会话 hash
func (s *ChatCompletionsService) GenerateSessionHash(req *openai.ChatCompletionsRequest) string {
	// 1. 从 user 字段提取 session_xxx
	if sessionID := openai.ExtractSessionID(req.User); sessionID != "" {
		return sessionID
	}

	// 2. 从 system 消息生成 hash
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			content := msg.GetContentAsString()
			if content != "" {
				return s.hashContent(content)
			}
		}
	}

	// 3. 从第一条 user 消息生成 hash
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			content := msg.GetContentAsString()
			if content != "" {
				return s.hashContent(content)
			}
			break
		}
	}

	return ""
}

func (s *ChatCompletionsService) hashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:16])
}

// SelectAccount 选择账户（复用 GatewayService 的逻辑）
func (s *ChatCompletionsService) SelectAccount(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	model string,
	failedAccountIDs map[int64]struct{},
) (*AccountSelectionResult, error) {
	// 复用 GatewayService 的账户选择逻辑
	return s.gatewayService.SelectAccountWithLoadAwareness(ctx, groupID, sessionHash, model, failedAccountIDs)
}

// BindStickySession 绑定粘性会话
func (s *ChatCompletionsService) BindStickySession(ctx context.Context, sessionHash string, accountID int64) error {
	return s.gatewayService.BindStickySession(ctx, sessionHash, accountID)
}

// Forward 转发 Chat Completions 请求
func (s *ChatCompletionsService) Forward(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	req *openai.ChatCompletionsRequest,
) (*ChatForwardResult, error) {
	startTime := time.Now()

	// 根据账户平台决定转发方式
	switch account.Platform {
	case PlatformAnthropic:
		return s.forwardToClaude(ctx, c, account, req, startTime)
	case PlatformAntigravity:
		// Antigravity 账户根据模型判断
		if isGeminiModel(req.Model) {
			return s.forwardToGemini(ctx, c, account, req, startTime)
		}
		return s.forwardToClaude(ctx, c, account, req, startTime)
	case PlatformGemini:
		return s.forwardToGemini(ctx, c, account, req, startTime)
	case PlatformOpenAI:
		return s.forwardToOpenAI(ctx, c, account, req, startTime)
	default:
		// 默认尝试 Claude 格式
		return s.forwardToClaude(ctx, c, account, req, startTime)
	}
}

// forwardToClaude 转发到 Claude API
func (s *ChatCompletionsService) forwardToClaude(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	req *openai.ChatCompletionsRequest,
	startTime time.Time,
) (*ChatForwardResult, error) {
	// 1. 转换请求
	claudeReq, err := openai.TransformChatToClaude(req, true) // 启用缓存
	if err != nil {
		return nil, fmt.Errorf("transform request: %w", err)
	}

	// 2. 获取映射的模型名称
	mappedModel := account.GetMappedModel(req.Model)
	if mappedModel == "" {
		mappedModel = req.Model
	}
	claudeReq.Model = mappedModel

	// 3. 序列化请求
	body, err := openai.SerializeClaudeRequest(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("serialize request: %w", err)
	}

	// 4. 获取凭证
	token, tokenType, err := s.gatewayService.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}

	// 4.5 OAuth 账户：注入 metadata.user_id 并重写
	if tokenType == "oauth" && s.identityService != nil {
		fp, fpErr := s.identityService.GetOrCreateFingerprint(ctx, account.ID, c.Request.Header)
		if fpErr == nil && fp != nil {
			accountUUID := account.GetExtraString("account_uuid")
			log.Printf("[ChatCompletions] OAuth account: clientID=%s, accountUUID=%s", fp.ClientID, accountUUID)
			if accountUUID != "" && fp.ClientID != "" {
				// 先注入 metadata.user_id（如果不存在）
				bodyBefore := string(body)
				body = s.injectMetadataUserID(body, fp.ClientID, accountUUID)
				log.Printf("[ChatCompletions] After injectMetadataUserID: %s", s.extractMetadataUserID(body))
				// 然后重写 user_id
				if newBody, err := s.identityService.RewriteUserID(body, account.ID, accountUUID, fp.ClientID); err == nil && len(newBody) > 0 {
					body = newBody
					log.Printf("[ChatCompletions] After RewriteUserID: %s", s.extractMetadataUserID(body))
				} else {
					log.Printf("[ChatCompletions] RewriteUserID skipped or failed, body unchanged from: %s", bodyBefore[:min(200, len(bodyBefore))])
				}
			} else {
				log.Printf("[ChatCompletions] Missing accountUUID or clientID, skipping metadata injection")
			}
		} else {
			log.Printf("[ChatCompletions] Failed to get fingerprint: %v", fpErr)
		}
	}



	// 5. 构建上游请求
	upstreamURL := "https://api.anthropic.com/v1/messages?beta=true"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// 根据凭证类型设置认证头和 beta header
	if tokenType == "oauth" {
		httpReq.Header.Set("authorization", "Bearer "+token)
		// OAuth 账户需要完整的 Claude Code beta headers
		// 根据模型决定使用哪个 header（Haiku 模型不需要 claude-code beta）
		betaHeader := claude.DefaultBetaHeader
		if strings.Contains(strings.ToLower(req.Model), "haiku") {
			betaHeader = claude.HaikuBetaHeader
		}
		httpReq.Header.Set("anthropic-beta", betaHeader)
		// 应用指纹 (Claude Code 客户端标识)
		if s.identityService != nil {
			fp, err := s.identityService.GetOrCreateFingerprint(ctx, account.ID, c.Request.Header)
			if err == nil && fp != nil {
				s.identityService.ApplyFingerprint(httpReq, fp)
			} else {
				// 降级：使用默认 Claude Code headers
				httpReq.Header.Set("User-Agent", "claude-cli/2.0.62 (external, cli)")
				httpReq.Header.Set("X-Stainless-Lang", "js")
				httpReq.Header.Set("X-Stainless-Package-Version", "0.52.0")
				httpReq.Header.Set("X-Stainless-OS", "Linux")
				httpReq.Header.Set("X-Stainless-Arch", "x64")
				httpReq.Header.Set("X-Stainless-Runtime", "node")
				httpReq.Header.Set("X-Stainless-Runtime-Version", "v22.14.0")
			}
		}
		httpReq.Header.Set("X-App", "cli")
		httpReq.Header.Set("Anthropic-Dangerous-Direct-Browser-Access", "true")
	} else {
		httpReq.Header.Set("x-api-key", token)
		httpReq.Header.Set("anthropic-beta", claude.APIKeyBetaHeader)
	}

	// 6. 发送请求
	resp, err := s.httpUpstream.Do(httpReq, "", account.ID, account.Concurrency)
	if err != nil {
		return nil, fmt.Errorf("upstream request: %w", err)
	}
	defer resp.Body.Close()

	// 7. 检查错误状态码
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("[ChatCompletions] Upstream error %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode}
	}

	// 8. 处理响应
	result := &ChatForwardResult{
		Model:  req.Model,
		Stream: req.Stream,
	}

	if req.Stream {
		err = s.handleClaudeStreamingResponse(c, resp, req.Model, req.StreamOptions, result, startTime)
	} else {
		err = s.handleClaudeNonStreamingResponse(c, resp, req.Model, result)
	}

	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// handleClaudeStreamingResponse 处理 Claude 流式响应
func (s *ChatCompletionsService) handleClaudeStreamingResponse(
	c *gin.Context,
	resp *http.Response,
	originalModel string,
	streamOptions *openai.StreamOptions,
	result *ChatForwardResult,
	startTime time.Time,
) error {
	includeUsage := streamOptions != nil && streamOptions.IncludeUsage

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	processor := openai.NewChatStreamingProcessor(originalModel, includeUsage)
	reader := bufio.NewReader(resp.Body)
	flusher, _ := c.Writer.(http.Flusher)
	firstTokenSent := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		output := processor.ProcessClaudeSSE(line)
		if len(output) > 0 {
			if !firstTokenSent {
				firstTokenSent = true
				ms := int(time.Since(startTime).Milliseconds())
				result.FirstTokenMs = &ms
			}
			if _, err := c.Writer.Write(output); err != nil {
				return err
			}
			if flusher != nil {
				flusher.Flush()
			}
		}
	}

	// 发送最终事件
	finalOutput := processor.Finish()
	if len(finalOutput) > 0 {
		if _, err := c.Writer.Write(finalOutput); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	result.Usage = processor.GetUsage()
	return nil
}

// handleClaudeNonStreamingResponse 处理 Claude 非流式响应
func (s *ChatCompletionsService) handleClaudeNonStreamingResponse(
	c *gin.Context,
	resp *http.Response,
	originalModel string,
	result *ChatForwardResult,
) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// 解析 Claude 响应
	var claudeResp openai.ClaudeResponseForTransform
	if err := json.Unmarshal(bodyBytes, &claudeResp); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	result.RequestID = claudeResp.ID

	// 转换为 OpenAI 格式
	chatResp := openai.TransformClaudeResponseToChat(&claudeResp, originalModel)

	result.Usage = chatResp.Usage

	// 发送响应
	c.JSON(http.StatusOK, chatResp)
	return nil
}

// forwardToGemini 转发到 Gemini API
func (s *ChatCompletionsService) forwardToGemini(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	req *openai.ChatCompletionsRequest,
	startTime time.Time,
) (*ChatForwardResult, error) {
	// TODO: 实现 Gemini 转发
	// 暂时返回错误
	return nil, fmt.Errorf("gemini forwarding not implemented yet")
}

// forwardToOpenAI 转发到 OpenAI API (原生格式)
func (s *ChatCompletionsService) forwardToOpenAI(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	req *openai.ChatCompletionsRequest,
	startTime time.Time,
) (*ChatForwardResult, error) {
	// 1. 序列化请求（直接使用 OpenAI 格式）
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("serialize request: %w", err)
	}

	// 2. 获取凭证
	token, _, err := s.gatewayService.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}

	// 3. 构建上游请求
	upstreamURL := "https://api.openai.com/v1/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	// 4. 发送请求
	resp, err := s.httpUpstream.Do(httpReq, "", account.ID, account.Concurrency)
	if err != nil {
		return nil, fmt.Errorf("upstream request: %w", err)
	}
	defer resp.Body.Close()

	// 5. 检查错误状态码
	if resp.StatusCode >= 400 {
		return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode}
	}

	// 6. 直接转发响应（OpenAI -> OpenAI 不需要转换）
	result := &ChatForwardResult{
		Model:  req.Model,
		Stream: req.Stream,
	}

	if req.Stream {
		err = s.proxyStreamingResponse(c, resp, result, startTime)
	} else {
		err = s.proxyNonStreamingResponse(c, resp, result)
	}

	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

// proxyStreamingResponse 代理流式响应（OpenAI 原生）
func (s *ChatCompletionsService) proxyStreamingResponse(
	c *gin.Context,
	resp *http.Response,
	result *ChatForwardResult,
	startTime time.Time,
) error {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	reader := bufio.NewReader(resp.Body)
	flusher, _ := c.Writer.(http.Flusher)
	firstTokenSent := false

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if !firstTokenSent && len(line) > 0 {
			firstTokenSent = true
			ms := int(time.Since(startTime).Milliseconds())
			result.FirstTokenMs = &ms
		}

		if _, err := c.Writer.Write(line); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	return nil
}

// proxyNonStreamingResponse 代理非流式响应（OpenAI 原生）
func (s *ChatCompletionsService) proxyNonStreamingResponse(
	c *gin.Context,
	resp *http.Response,
	result *ChatForwardResult,
) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// 解析用量信息
	var chatResp openai.ChatCompletionsResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err == nil {
		result.RequestID = chatResp.ID
		result.Usage = chatResp.Usage
	}

	c.Header("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, _ = c.Writer.Write(bodyBytes)
	return nil
}

// RecordUsage 记录使用量
func (s *ChatCompletionsService) RecordUsage(
	ctx context.Context,
	result *ChatForwardResult,
	apiKey *APIKey,
	user *User,
	account *Account,
	subscription *UserSubscription,
) error {
	if result.Usage == nil {
		return nil
	}

	// 使用 BillingService 计算费用
	tokens := UsageTokens{
		InputTokens:  result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
	}

	// 获取缓存 token
	if result.Usage.PromptTokensDetails != nil {
		tokens.CacheReadTokens = result.Usage.PromptTokensDetails.CachedTokens
	}

	// 计算费用
	var cost float64
	rateMultiplier := 1.0
	if apiKey.Group != nil {
		rateMultiplier = apiKey.Group.RateMultiplier
	}

	breakdown, err := s.billingService.CalculateCost(result.Model, tokens, rateMultiplier)
	if err != nil {
		log.Printf("[ChatCompletions] Calculate cost error: %v", err)
	} else {
		cost = breakdown.ActualCost
	}

	// 记录用量日志（复用现有逻辑）
	log.Printf("[ChatCompletions] Usage recorded: model=%s input=%d output=%d cost=%.6f",
		result.Model, tokens.InputTokens, tokens.OutputTokens, cost)

	return nil
}

// isGeminiModel 判断是否是 Gemini 模型
func isGeminiModel(model string) bool {
	model = strings.ToLower(model)
	return strings.HasPrefix(model, "gemini") || strings.Contains(model, "gemini")
}

// injectMetadataUserID 为 Claude 请求注入 metadata.user_id
// 格式必须匹配 RewriteUserID 期望的正则: user_{64位hex}_account__session_{uuid}
func (s *ChatCompletionsService) injectMetadataUserID(body []byte, clientID, accountUUID string) []byte {
	if len(body) == 0 || clientID == "" {
		return body
	}

	var reqMap map[string]any
	if err := json.Unmarshal(body, &reqMap); err != nil {
		return body
	}

	// 检查是否已有 metadata
	metadata, ok := reqMap["metadata"].(map[string]any)
	if !ok {
		metadata = make(map[string]any)
	}

	// 检查是否已有 user_id
	if _, hasUserID := metadata["user_id"]; !hasUserID {
		// 生成 Claude Code 风格的 user_id
		// 格式必须是: user_{64位hex}_account__session_{uuid}
		// 注意: _account__session_ 是双下划线，中间没有 accountUUID

		// 生成 36 字符 UUID 格式的 session ID
		sessionUUID := s.generateUUID()

		// clientID 必须是 64 位 hex，如果不够则补齐
		paddedClientID := clientID
		if len(paddedClientID) < 64 {
			// 用 hash 补齐到 64 位
			h := sha256.Sum256([]byte(clientID))
			paddedClientID = hex.EncodeToString(h[:])
		} else if len(paddedClientID) > 64 {
			paddedClientID = paddedClientID[:64]
		}

		// 格式: user_{64位hex}_account__session_{uuid}
		metadata["user_id"] = fmt.Sprintf("user_%s_account__session_%s", paddedClientID, sessionUUID)
	}

	reqMap["metadata"] = metadata

	newBody, err := json.Marshal(reqMap)
	if err != nil {
		return body
	}
	return newBody
}

// generateUUID 生成标准 UUID 格式字符串 (8-4-4-4-12)
func (s *ChatCompletionsService) generateUUID() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	hexStr := hex.EncodeToString(h[:18])
	return fmt.Sprintf("%s-%s-%s-%s-%s", hexStr[0:8], hexStr[8:12], hexStr[12:16], hexStr[16:20], hexStr[20:32])
}

// extractMetadataUserID 从请求体中提取 metadata.user_id（用于调试日志）
func (s *ChatCompletionsService) extractMetadataUserID(body []byte) string {
	var reqMap map[string]any
	if err := json.Unmarshal(body, &reqMap); err != nil {
		return "(parse error)"
	}
	metadata, ok := reqMap["metadata"].(map[string]any)
	if !ok {
		return "(no metadata)"
	}
	userID, ok := metadata["user_id"].(string)
	if !ok {
		return "(no user_id)"
	}
	return userID
}
