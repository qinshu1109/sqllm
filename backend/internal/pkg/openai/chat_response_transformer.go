package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ============================================================
// Claude -> OpenAI Chat Completions 响应转换
// ============================================================

// ClaudeResponseForTransform Claude 响应（用于转换）
type ClaudeResponseForTransform struct {
	ID           string                    `json:"id"`
	Type         string                    `json:"type"`
	Role         string                    `json:"role"`
	Model        string                    `json:"model"`
	Content      []ClaudeResponseContent   `json:"content"`
	StopReason   string                    `json:"stop_reason,omitempty"`
	StopSequence *string                   `json:"stop_sequence,omitempty"`
	Usage        ClaudeUsageForTransform   `json:"usage"`
}

// ClaudeResponseContent Claude 响应内容项
type ClaudeResponseContent struct {
	Type      string `json:"type"` // text, thinking, tool_use
	Text      string `json:"text,omitempty"`
	Thinking  string `json:"thinking,omitempty"`
	Signature string `json:"signature,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Input     any    `json:"input,omitempty"`
}

// ClaudeUsageForTransform Claude 用量统计
type ClaudeUsageForTransform struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

// TransformClaudeResponseToChat 将 Claude 响应转换为 OpenAI Chat Completions 格式
func TransformClaudeResponseToChat(claudeResp *ClaudeResponseForTransform, originalModel string) *ChatCompletionsResponse {
	requestID := claudeResp.ID
	if requestID == "" {
		requestID = "chatcmpl-" + generateID()
	}

	model := originalModel
	if model == "" {
		model = claudeResp.Model
	}

	resp := &ChatCompletionsResponse{
		ID:      requestID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatChoice{
			{
				Index:        0,
				Message:      transformClaudeContentToMessage(claudeResp.Content),
				FinishReason: mapClaudeStopReason(claudeResp.StopReason),
			},
		},
		Usage: &ChatUsage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		},
	}

	// 添加缓存 token 详情
	if claudeResp.Usage.CacheReadInputTokens > 0 {
		resp.Usage.PromptTokensDetails = &PromptTokensDetails{
			CachedTokens: claudeResp.Usage.CacheReadInputTokens,
		}
	}

	return resp
}

// transformClaudeContentToMessage 转换 Claude 内容为 OpenAI 消息
func transformClaudeContentToMessage(content []ClaudeResponseContent) *ChatMessage {
	msg := &ChatMessage{
		Role: "assistant",
	}

	var textParts []string
	var toolCalls []ToolCall

	for _, item := range content {
		switch item.Type {
		case "text":
			textParts = append(textParts, item.Text)
		case "thinking":
			// 忽略 thinking 块，或可选择性包含
			// textParts = append(textParts, fmt.Sprintf("<thinking>%s</thinking>", item.Thinking))
		case "tool_use":
			inputJSON, _ := json.Marshal(item.Input)
			toolCalls = append(toolCalls, ToolCall{
				ID:   item.ID,
				Type: "function",
				Function: FunctionCall{
					Name:      item.Name,
					Arguments: string(inputJSON),
				},
			})
		}
	}

	if len(textParts) > 0 {
		msg.Content = strings.Join(textParts, "")
	}
	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}

	return msg
}

// mapClaudeStopReason 映射 Claude stop_reason 到 OpenAI finish_reason
func mapClaudeStopReason(reason string) *string {
	var mapped string
	switch reason {
	case "end_turn":
		mapped = "stop"
	case "tool_use":
		mapped = "tool_calls"
	case "max_tokens":
		mapped = "length"
	case "stop_sequence":
		mapped = "stop"
	default:
		if reason != "" {
			mapped = "stop"
		} else {
			return nil
		}
	}
	return &mapped
}

// ============================================================
// Streaming Processor - Claude SSE -> OpenAI SSE
// ============================================================

// ChatStreamingProcessor 流式响应处理器
type ChatStreamingProcessor struct {
	requestID        string
	originalModel    string
	created          int64
	messageStartSent bool
	includeUsage     bool

	// 当前状态
	currentTextIndex int
	toolCallsStarted bool
	toolCallIndex    int

	// 累计 usage
	inputTokens      int
	outputTokens     int
	cacheReadTokens  int
}

// NewChatStreamingProcessor 创建流式响应处理器
func NewChatStreamingProcessor(originalModel string, includeUsage bool) *ChatStreamingProcessor {
	return &ChatStreamingProcessor{
		requestID:        "chatcmpl-" + generateID(),
		originalModel:    originalModel,
		created:          time.Now().Unix(),
		includeUsage:     includeUsage,
		currentTextIndex: -1,
		toolCallIndex:    -1,
	}
}

// ProcessClaudeSSE 处理 Claude SSE 事件，返回 OpenAI SSE
func (p *ChatStreamingProcessor) ProcessClaudeSSE(line string) []byte {
	line = strings.TrimSpace(line)

	// 处理 event: 行
	if strings.HasPrefix(line, "event:") {
		return nil // 只处理 data: 行
	}

	if !strings.HasPrefix(line, "data:") {
		return nil
	}

	data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if data == "" || data == "[DONE]" {
		return nil
	}

	// 解析 Claude SSE 事件
	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return nil
	}

	eventType, _ := event["type"].(string)

	var result bytes.Buffer

	switch eventType {
	case "message_start":
		result.Write(p.handleMessageStart(event))

	case "content_block_start":
		result.Write(p.handleContentBlockStart(event))

	case "content_block_delta":
		result.Write(p.handleContentBlockDelta(event))

	case "content_block_stop":
		// 不需要特殊处理

	case "message_delta":
		result.Write(p.handleMessageDelta(event))

	case "message_stop":
		// 返回 [DONE]
		result.WriteString("data: [DONE]\n\n")
	}

	return result.Bytes()
}

// handleMessageStart 处理 message_start 事件
func (p *ChatStreamingProcessor) handleMessageStart(event map[string]any) []byte {
	if p.messageStartSent {
		return nil
	}
	p.messageStartSent = true

	// 提取 usage
	if message, ok := event["message"].(map[string]any); ok {
		if usage, ok := message["usage"].(map[string]any); ok {
			if v, ok := usage["input_tokens"].(float64); ok {
				p.inputTokens = int(v)
			}
			if v, ok := usage["cache_read_input_tokens"].(float64); ok {
				p.cacheReadTokens = int(v)
			}
		}
	}

	// 发送第一个 chunk（带 role）
	chunk := p.buildChunk(&ChatDelta{Role: "assistant"}, nil)
	return p.formatSSE(chunk)
}

// handleContentBlockStart 处理 content_block_start 事件
func (p *ChatStreamingProcessor) handleContentBlockStart(event map[string]any) []byte {
	contentBlock, ok := event["content_block"].(map[string]any)
	if !ok {
		return nil
	}

	blockType, _ := contentBlock["type"].(string)

	switch blockType {
	case "text":
		p.currentTextIndex++
		// 不需要发送特殊事件

	case "tool_use":
		p.toolCallsStarted = true
		p.toolCallIndex++

		id, _ := contentBlock["id"].(string)
		name, _ := contentBlock["name"].(string)

		// 发送工具调用开始
		chunk := p.buildChunk(&ChatDelta{
			ToolCalls: []ToolCall{{
				Index: &p.toolCallIndex,
				ID:    id,
				Type:  "function",
				Function: FunctionCall{
					Name:      name,
					Arguments: "",
				},
			}},
		}, nil)
		return p.formatSSE(chunk)
	}

	return nil
}

// handleContentBlockDelta 处理 content_block_delta 事件
func (p *ChatStreamingProcessor) handleContentBlockDelta(event map[string]any) []byte {
	delta, ok := event["delta"].(map[string]any)
	if !ok {
		return nil
	}

	deltaType, _ := delta["type"].(string)

	switch deltaType {
	case "text_delta":
		text, _ := delta["text"].(string)
		if text == "" {
			return nil
		}
		chunk := p.buildChunk(&ChatDelta{Content: text}, nil)
		return p.formatSSE(chunk)

	case "input_json_delta":
		partialJSON, _ := delta["partial_json"].(string)
		if partialJSON == "" {
			return nil
		}
		// 发送工具调用参数增量
		idx := p.toolCallIndex
		chunk := p.buildChunk(&ChatDelta{
			ToolCalls: []ToolCall{{
				Index: &idx,
				Function: FunctionCall{
					Arguments: partialJSON,
				},
			}},
		}, nil)
		return p.formatSSE(chunk)

	case "thinking_delta":
		// 忽略 thinking 增量
		return nil
	}

	return nil
}

// handleMessageDelta 处理 message_delta 事件
func (p *ChatStreamingProcessor) handleMessageDelta(event map[string]any) []byte {
	delta, ok := event["delta"].(map[string]any)
	if !ok {
		return nil
	}

	// 提取 stop_reason
	stopReason, _ := delta["stop_reason"].(string)
	finishReason := mapClaudeStopReason(stopReason)

	// 提取 usage
	if usage, ok := event["usage"].(map[string]any); ok {
		if v, ok := usage["output_tokens"].(float64); ok {
			p.outputTokens = int(v)
		}
	}

	// 发送最终 chunk
	chunk := p.buildChunk(&ChatDelta{}, finishReason)
	return p.formatSSE(chunk)
}

// buildChunk 构建 OpenAI chunk
func (p *ChatStreamingProcessor) buildChunk(delta *ChatDelta, finishReason *string) *ChatCompletionsChunk {
	return &ChatCompletionsChunk{
		ID:      p.requestID,
		Object:  "chat.completion.chunk",
		Created: p.created,
		Model:   p.originalModel,
		Choices: []ChatChoice{{
			Index:        0,
			Delta:        delta,
			FinishReason: finishReason,
		}},
	}
}

// formatSSE 格式化为 SSE
func (p *ChatStreamingProcessor) formatSSE(chunk *ChatCompletionsChunk) []byte {
	data, err := json.Marshal(chunk)
	if err != nil {
		return nil
	}
	return []byte(fmt.Sprintf("data: %s\n\n", data))
}

// Finish 结束处理，返回最终 usage 事件
func (p *ChatStreamingProcessor) Finish() []byte {
	if !p.includeUsage {
		return []byte("data: [DONE]\n\n")
	}

	// 发送带 usage 的最终 chunk
	chunk := &ChatCompletionsChunk{
		ID:      p.requestID,
		Object:  "chat.completion.chunk",
		Created: p.created,
		Model:   p.originalModel,
		Choices: []ChatChoice{},
		Usage: &ChatUsage{
			PromptTokens:     p.inputTokens,
			CompletionTokens: p.outputTokens,
			TotalTokens:      p.inputTokens + p.outputTokens,
		},
	}

	if p.cacheReadTokens > 0 {
		chunk.Usage.PromptTokensDetails = &PromptTokensDetails{
			CachedTokens: p.cacheReadTokens,
		}
	}

	var result bytes.Buffer
	result.Write(p.formatSSE(chunk))
	result.WriteString("data: [DONE]\n\n")
	return result.Bytes()
}

// GetUsage 获取累计用量
func (p *ChatStreamingProcessor) GetUsage() *ChatUsage {
	return &ChatUsage{
		PromptTokens:     p.inputTokens,
		CompletionTokens: p.outputTokens,
		TotalTokens:      p.inputTokens + p.outputTokens,
		PromptTokensDetails: func() *PromptTokensDetails {
			if p.cacheReadTokens > 0 {
				return &PromptTokensDetails{CachedTokens: p.cacheReadTokens}
			}
			return nil
		}(),
	}
}

// ============================================================
// Helper Functions
// ============================================================

func generateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:24]
}

// BuildErrorResponse 构建 OpenAI 错误响应
func BuildErrorResponse(errType, message string) *ChatErrorResponse {
	return &ChatErrorResponse{
		Error: ChatError{
			Type:    errType,
			Message: message,
		},
	}
}

// FormatErrorSSE 格式化错误为 SSE
func FormatErrorSSE(errType, message string) []byte {
	resp := BuildErrorResponse(errType, message)
	data, _ := json.Marshal(resp)
	return []byte(fmt.Sprintf("data: %s\n\ndata: [DONE]\n\n", data))
}
