// Package service provides OpenAI Chat Completions proxy with Claude caching support.
package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	// OpenAI Chat 代理的系统提示词（与 Claude 原生 API 保持一致）
	openaiChatClaudeCodePrompt = "You are Claude Code, Anthropic's official CLI for Claude."
	// 最大 cache_control 块数量
	openaiChatMaxCacheBlocks = 4
)

// OpenAIChatUsage represents usage statistics from LiteLLM response
type OpenAIChatUsage struct {
	PromptTokens            int `json:"prompt_tokens"`
	CompletionTokens        int `json:"completion_tokens"`
	TotalTokens             int `json:"total_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens    int `json:"cache_read_input_tokens,omitempty"`
}

// OpenAIChatForwardResult represents the result of forwarding to LiteLLM
type OpenAIChatForwardResult struct {
	RequestID string
	Model     string
	Stream    bool
	Usage     OpenAIChatUsage
	Duration  time.Duration
}

// OpenAIChatService handles OpenAI Chat Completions API with Claude caching
type OpenAIChatService struct {
	httpClient *http.Client
	cfg        *config.Config
}

// NewOpenAIChatService creates a new OpenAIChatService
func NewOpenAIChatService(cfg *config.Config) *OpenAIChatService {
	return &OpenAIChatService{
		httpClient: &http.Client{
			Timeout: 10 * time.Minute, // LLM 请求可能较长
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		cfg: cfg,
	}
}

// Forward forwards the request to LiteLLM with optional cache injection
func (s *OpenAIChatService) Forward(ctx context.Context, c *gin.Context, body []byte) (*OpenAIChatForwardResult, error) {
	startTime := time.Now()

	// 解析请求
	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, fmt.Errorf("parse request body: %w", err)
	}

	model, _ := reqBody["model"].(string)
	stream, _ := reqBody["stream"].(bool)

	// 判断是否需要注入缓存控制
	if s.shouldInjectCacheControl(model) {
		body = s.injectCacheControlForOpenAI(body, reqBody)
	}

	// 构建转发请求
	litellmURL := s.getLiteLLMURL()
	req, err := http.NewRequestWithContext(ctx, "POST", litellmURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 复制必要的请求头
	req.Header.Set("Content-Type", "application/json")
	if auth := c.GetHeader("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("forward request: %w", err)
	}

	// 处理响应
	if stream {
		return s.handleStreamingResponse(c, resp, model, startTime)
	}
	return s.handleNonStreamingResponse(c, resp, model, startTime)
}

// shouldInjectCacheControl determines if cache control should be injected
func (s *OpenAIChatService) shouldInjectCacheControl(model string) bool {
	// 检查配置是否启用
	if s.cfg != nil && !s.cfg.Gateway.LiteLLM.EnableClaudeCaching {
		return false
	}

	// 必须是 Claude 模型
	if !s.isClaudeModel(model) {
		return false
	}

	// Haiku 模型不支持缓存
	if strings.Contains(strings.ToLower(model), "haiku") {
		return false
	}

	return true
}

// isClaudeModel checks if the model is a Claude model
func (s *OpenAIChatService) isClaudeModel(model string) bool {
	lower := strings.ToLower(model)
	return strings.Contains(lower, "claude") ||
		strings.HasPrefix(lower, "anthropic/")
}

// injectCacheControlForOpenAI injects cache_control into OpenAI format request
func (s *OpenAIChatService) injectCacheControlForOpenAI(body []byte, reqBody map[string]any) []byte {
	messages, ok := reqBody["messages"].([]any)
	if !ok || len(messages) == 0 {
		return body
	}

	// 查找 system 消息
	systemIdx := -1
	var systemMsg map[string]any
	for i, msg := range messages {
		if m, ok := msg.(map[string]any); ok {
			if m["role"] == "system" {
				systemIdx = i
				systemMsg = m
				break
			}
		}
	}

	// 检查是否已包含 Claude Code 提示词
	if s.systemIncludesClaudeCodePromptOpenAI(messages) {
		return body
	}

	// 创建带 cache_control 的提示词块
	claudeCodeBlock := map[string]any{
		"type":          "text",
		"text":          openaiChatClaudeCodePrompt,
		"cache_control": map[string]string{"type": "ephemeral"},
	}

	if systemMsg == nil {
		// 没有 system 消息，创建一个新的
		newSystem := map[string]any{
			"role":    "system",
			"content": []any{claudeCodeBlock},
		}
		messages = append([]any{newSystem}, messages...)
	} else {
		// 有 system 消息，处理 content
		content := systemMsg["content"]
		var newContent []any

		switch v := content.(type) {
		case string:
			// 字符串格式转换为数组格式
			if v == "" || v == openaiChatClaudeCodePrompt {
				newContent = []any{claudeCodeBlock}
			} else {
				newContent = []any{
					claudeCodeBlock,
					map[string]any{"type": "text", "text": v},
				}
			}
		case []any:
			// 已经是数组格式，在开头插入
			newContent = make([]any, 0, len(v)+1)
			newContent = append(newContent, claudeCodeBlock)
			for _, item := range v {
				// 过滤掉重复的 Claude Code 提示词
				if m, ok := item.(map[string]any); ok {
					if text, ok := m["text"].(string); ok && text == openaiChatClaudeCodePrompt {
						continue
					}
				}
				newContent = append(newContent, item)
			}
		default:
			newContent = []any{claudeCodeBlock}
		}

		systemMsg["content"] = newContent
		messages[systemIdx] = systemMsg
	}

	reqBody["messages"] = messages

	// 强制执行 cache_control 块数量限制
	s.enforceCacheControlLimitOpenAI(reqBody)

	// 序列化
	result, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Warning: failed to marshal request after cache injection: %v", err)
		return body
	}

	return result
}

// systemIncludesClaudeCodePromptOpenAI checks if messages already include Claude Code prompt
func (s *OpenAIChatService) systemIncludesClaudeCodePromptOpenAI(messages []any) bool {
	for _, msg := range messages {
		m, ok := msg.(map[string]any)
		if !ok || m["role"] != "system" {
			continue
		}

		content := m["content"]
		switch v := content.(type) {
		case string:
			if s.hasClaudeCodePrefix(v) {
				return true
			}
		case []any:
			for _, item := range v {
				if itemMap, ok := item.(map[string]any); ok {
					if text, ok := itemMap["text"].(string); ok && s.hasClaudeCodePrefix(text) {
						return true
					}
				}
			}
		}
	}
	return false
}

// hasClaudeCodePrefix checks if text starts with Claude Code prompt prefix
func (s *OpenAIChatService) hasClaudeCodePrefix(text string) bool {
	prefixes := []string{
		"You are Claude Code, Anthropic's official CLI for Claude",
		"You are a Claude agent, built on Anthropic's Claude Agent SDK",
		"You are a file search specialist for Claude Code",
		"You are a helpful AI assistant tasked with summarizing conversations",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(text, prefix) {
			return true
		}
	}
	return false
}

// enforceCacheControlLimitOpenAI enforces the max cache_control blocks limit (4)
func (s *OpenAIChatService) enforceCacheControlLimitOpenAI(reqBody map[string]any) {
	count := s.countCacheControlBlocksOpenAI(reqBody)
	if count <= openaiChatMaxCacheBlocks {
		return
	}

	messages, ok := reqBody["messages"].([]any)
	if !ok {
		return
	}

	// 超限时从非 system 消息中移除 cache_control
	for count > openaiChatMaxCacheBlocks {
		removed := false
		for _, msg := range messages {
			m, ok := msg.(map[string]any)
			if !ok || m["role"] == "system" {
				continue
			}
			if content, ok := m["content"].([]any); ok {
				for _, item := range content {
					if itemMap, ok := item.(map[string]any); ok {
						if _, has := itemMap["cache_control"]; has {
							delete(itemMap, "cache_control")
							count--
							removed = true
							break
						}
					}
				}
			}
			if removed {
				break
			}
		}
		if !removed {
			break
		}
	}
}

// countCacheControlBlocksOpenAI counts cache_control blocks in OpenAI format
func (s *OpenAIChatService) countCacheControlBlocksOpenAI(reqBody map[string]any) int {
	count := 0
	messages, ok := reqBody["messages"].([]any)
	if !ok {
		return 0
	}

	for _, msg := range messages {
		m, ok := msg.(map[string]any)
		if !ok {
			continue
		}
		if content, ok := m["content"].([]any); ok {
			for _, item := range content {
				if itemMap, ok := item.(map[string]any); ok {
					if _, has := itemMap["cache_control"]; has {
						count++
					}
				}
			}
		}
	}
	return count
}

// getLiteLLMURL returns the LiteLLM backend URL
func (s *OpenAIChatService) getLiteLLMURL() string {
	if s.cfg != nil && s.cfg.Gateway.LiteLLM.URL != "" {
		return strings.TrimSuffix(s.cfg.Gateway.LiteLLM.URL, "/")
	}
	return "http://localhost:4000"
}

// handleStreamingResponse handles SSE streaming response from LiteLLM
func (s *OpenAIChatService) handleStreamingResponse(c *gin.Context, resp *http.Response, model string, startTime time.Time) (*OpenAIChatForwardResult, error) {
	defer resp.Body.Close()

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 透传响应
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	var usage OpenAIChatUsage
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024) // 支持大行

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(c.Writer, "%s\n", line)
		flusher.Flush()

		// 尝试从 SSE 数据中解析 usage
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "[DONE]" {
				var chunk struct {
					Usage *OpenAIChatUsage `json:"usage,omitempty"`
				}
				if err := json.Unmarshal([]byte(data), &chunk); err == nil && chunk.Usage != nil {
					usage = *chunk.Usage
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read stream: %w", err)
	}

	return &OpenAIChatForwardResult{
		Model:    model,
		Stream:   true,
		Usage:    usage,
		Duration: time.Since(startTime),
	}, nil
}

// handleNonStreamingResponse handles non-streaming response from LiteLLM
func (s *OpenAIChatService) handleNonStreamingResponse(c *gin.Context, resp *http.Response, model string, startTime time.Time) (*OpenAIChatForwardResult, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// 设置响应头并返回响应体
	c.Header("Content-Type", "application/json")
	c.Status(resp.StatusCode)
	c.Writer.Write(body)

	// 解析 usage
	var respBody struct {
		Usage OpenAIChatUsage `json:"usage"`
	}
	_ = json.Unmarshal(body, &respBody)

	return &OpenAIChatForwardResult{
		Model:    model,
		Stream:   false,
		Usage:    respBody.Usage,
		Duration: time.Since(startTime),
	}, nil
}
