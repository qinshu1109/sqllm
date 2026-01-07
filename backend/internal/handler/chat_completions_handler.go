package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ChatCompletionsHandler 处理 OpenAI Chat Completions API 请求
type ChatCompletionsHandler struct {
	chatService         *service.ChatCompletionsService
	billingCacheService *service.BillingCacheService
	concurrencyHelper   *ConcurrencyHelper
}

// NewChatCompletionsHandler 创建 ChatCompletionsHandler
func NewChatCompletionsHandler(
	chatService *service.ChatCompletionsService,
	concurrencyService *service.ConcurrencyService,
	billingCacheService *service.BillingCacheService,
	cfg *config.Config,
) *ChatCompletionsHandler {
	pingInterval := time.Duration(0)
	if cfg != nil {
		pingInterval = time.Duration(cfg.Concurrency.PingInterval) * time.Second
	}
	return &ChatCompletionsHandler{
		chatService:         chatService,
		billingCacheService: billingCacheService,
		concurrencyHelper:   NewConcurrencyHelper(concurrencyService, SSEPingFormatOpenAI, pingInterval),
	}
}

// SSEPingFormatOpenAI OpenAI 格式的 SSE ping
const SSEPingFormatOpenAI = ": ping\n\n"

// ChatCompletions 处理 POST /v1/chat/completions
func (h *ChatCompletionsHandler) ChatCompletions(c *gin.Context) {
	// 1. 从 context 获取 apiKey 和 user
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
	}

	// 2. 读取并解析请求体
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

	var req openai.ChatCompletionsRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body: "+err.Error())
		return
	}

	// 3. 验证必填字段
	if req.Model == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
	}

	if len(req.Messages) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "messages is required")
		return
	}

	// Track if we've started streaming
	streamStarted := false

	// 获取订阅信息
	subscription, _ := middleware2.GetSubscriptionFromContext(c)

	// 4. 检查 wait 队列
	maxWait := service.CalculateMaxWait(subject.Concurrency)
	canWait, err := h.concurrencyHelper.IncrementWaitCount(c.Request.Context(), subject.UserID, maxWait)
	if err != nil {
		log.Printf("Increment wait count failed: %v", err)
	} else if !canWait {
		h.errorResponse(c, http.StatusTooManyRequests, "rate_limit_error", "Too many pending requests, please retry later")
		return
	}
	defer h.concurrencyHelper.DecrementWaitCount(c.Request.Context(), subject.UserID)

	// 5. 获取用户并发槽位
	userReleaseFunc, err := h.concurrencyHelper.AcquireUserSlotWithWait(c, subject.UserID, subject.Concurrency, req.Stream, &streamStarted)
	if err != nil {
		log.Printf("User concurrency acquire failed: %v", err)
		h.handleConcurrencyError(c, err, "user", streamStarted)
		return
	}
	userReleaseFunc = wrapReleaseOnDone(c.Request.Context(), userReleaseFunc)
	if userReleaseFunc != nil {
		defer userReleaseFunc()
	}

	// 6. 二次检查余额/订阅
	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription); err != nil {
		log.Printf("Billing eligibility check failed: %v", err)
		status, code, message := billingErrorDetails(err)
		h.handleStreamingAwareError(c, status, code, message, streamStarted)
		return
	}

	// 7. 生成粘性会话 hash
	sessionHash := h.chatService.GenerateSessionHash(&req)

	// 8. 账户选择和转发
	const maxAccountSwitches = 3
	switchCount := 0
	failedAccountIDs := make(map[int64]struct{})
	lastFailoverStatus := 0

	for {
		// 选择账户
		log.Printf("[ChatCompletions] Selecting account: groupID=%v model=%s", apiKey.GroupID, req.Model)
		selection, err := h.chatService.SelectAccount(c.Request.Context(), apiKey.GroupID, sessionHash, req.Model, failedAccountIDs)
		if err != nil {
			log.Printf("[ChatCompletions] SelectAccount failed: %v", err)
			if len(failedAccountIDs) == 0 {
				h.handleStreamingAwareError(c, http.StatusServiceUnavailable, "api_error", "No available accounts: "+err.Error(), streamStarted)
				return
			}
			h.handleFailoverExhausted(c, lastFailoverStatus, streamStarted)
			return
		}
		account := selection.Account
		log.Printf("[ChatCompletions] Selected account: id=%d name=%s platform=%s", account.ID, account.Name, account.Platform)

		// 获取账户并发槽位
		accountReleaseFunc := selection.ReleaseFunc
		var accountWaitRelease func()
		if !selection.Acquired {
			if selection.WaitPlan == nil {
				h.handleStreamingAwareError(c, http.StatusServiceUnavailable, "api_error", "No available accounts", streamStarted)
				return
			}
			canWait, err := h.concurrencyHelper.IncrementAccountWaitCount(c.Request.Context(), account.ID, selection.WaitPlan.MaxWaiting)
			if err != nil {
				log.Printf("Increment account wait count failed: %v", err)
			} else if !canWait {
				h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error", "Too many pending requests, please retry later", streamStarted)
				return
			} else {
				accountWaitRelease = func() {
					h.concurrencyHelper.DecrementAccountWaitCount(c.Request.Context(), account.ID)
				}
			}

			accountReleaseFunc, err = h.concurrencyHelper.AcquireAccountSlotWithWaitTimeout(
				c,
				account.ID,
				selection.WaitPlan.MaxConcurrency,
				selection.WaitPlan.Timeout,
				req.Stream,
				&streamStarted,
			)
			if err != nil {
				if accountWaitRelease != nil {
					accountWaitRelease()
				}
				log.Printf("Account concurrency acquire failed: %v", err)
				h.handleConcurrencyError(c, err, "account", streamStarted)
				return
			}
			if err := h.chatService.BindStickySession(c.Request.Context(), sessionHash, account.ID); err != nil {
				log.Printf("Bind sticky session failed: %v", err)
			}
		}
		accountReleaseFunc = wrapReleaseOnDone(c.Request.Context(), accountReleaseFunc)
		accountWaitRelease = wrapReleaseOnDone(c.Request.Context(), accountWaitRelease)

		// 转发请求
		result, err := h.chatService.Forward(c.Request.Context(), c, account, &req)
		if accountReleaseFunc != nil {
			accountReleaseFunc()
		}
		if accountWaitRelease != nil {
			accountWaitRelease()
		}

		if err != nil {
			var failoverErr *service.UpstreamFailoverError
			if errors.As(err, &failoverErr) {
				failedAccountIDs[account.ID] = struct{}{}
				if switchCount >= maxAccountSwitches {
					lastFailoverStatus = failoverErr.StatusCode
					h.handleFailoverExhausted(c, lastFailoverStatus, streamStarted)
					return
				}
				lastFailoverStatus = failoverErr.StatusCode
				switchCount++
				log.Printf("Account %d: upstream error %d, switching account %d/%d", account.ID, failoverErr.StatusCode, switchCount, maxAccountSwitches)
				continue
			}
			log.Printf("Account %d: Forward request failed: %v", account.ID, err)
			return
		}

		// 异步记录使用量
		go func(result *service.ChatForwardResult, usedAccount *service.Account) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := h.chatService.RecordUsage(ctx, result, apiKey, apiKey.User, usedAccount, subscription); err != nil {
				log.Printf("Record usage failed: %v", err)
			}
		}(result, account)
		return
	}
}

// handleConcurrencyError 处理并发错误
func (h *ChatCompletionsHandler) handleConcurrencyError(c *gin.Context, err error, slotType string, streamStarted bool) {
	h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error",
		fmt.Sprintf("Concurrency limit exceeded for %s, please retry later", slotType), streamStarted)
}

// handleFailoverExhausted 处理故障转移耗尽
func (h *ChatCompletionsHandler) handleFailoverExhausted(c *gin.Context, statusCode int, streamStarted bool) {
	status, errType, errMsg := h.mapUpstreamError(statusCode)
	h.handleStreamingAwareError(c, status, errType, errMsg, streamStarted)
}

// mapUpstreamError 映射上游错误
func (h *ChatCompletionsHandler) mapUpstreamError(statusCode int) (int, string, string) {
	switch statusCode {
	case 401:
		return http.StatusBadGateway, "upstream_error", "Upstream authentication failed, please contact administrator"
	case 403:
		return http.StatusBadGateway, "upstream_error", "Upstream access forbidden, please contact administrator"
	case 429:
		return http.StatusTooManyRequests, "rate_limit_error", "Upstream rate limit exceeded, please retry later"
	case 529:
		return http.StatusServiceUnavailable, "upstream_error", "Upstream service overloaded, please retry later"
	case 500, 502, 503, 504:
		return http.StatusBadGateway, "upstream_error", "Upstream service temporarily unavailable"
	default:
		return http.StatusBadGateway, "upstream_error", "Upstream request failed"
	}
}

// handleStreamingAwareError 处理流式感知错误
func (h *ChatCompletionsHandler) handleStreamingAwareError(c *gin.Context, status int, errType, message string, streamStarted bool) {
	if streamStarted {
		// 流已开始，发送 SSE 错误事件
		flusher, ok := c.Writer.(http.Flusher)
		if ok {
			errData := openai.FormatErrorSSE(errType, message)
			if _, err := c.Writer.Write(errData); err != nil {
				_ = c.Error(err)
			}
			flusher.Flush()
		}
		return
	}

	// 正常 JSON 错误响应
	h.errorResponse(c, status, errType, message)
}

// errorResponse 返回 OpenAI 格式错误响应
func (h *ChatCompletionsHandler) errorResponse(c *gin.Context, status int, errType, message string) {
	c.JSON(status, openai.BuildErrorResponse(errType, message))
}
