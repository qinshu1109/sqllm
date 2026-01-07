package openai

import "encoding/json"

// ============================================================
// OpenAI Chat Completions API Types
// https://platform.openai.com/docs/api-reference/chat/create
// ============================================================

// ChatCompletionsRequest OpenAI Chat Completions 请求
type ChatCompletionsRequest struct {
	Model               string          `json:"model"`
	Messages            []ChatMessage   `json:"messages"`
	Stream              bool            `json:"stream,omitempty"`
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	MaxTokens           *int            `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int            `json:"max_completion_tokens,omitempty"`
	Stop                any             `json:"stop,omitempty"` // string 或 []string
	Tools               []ChatTool      `json:"tools,omitempty"`
	ToolChoice          any             `json:"tool_choice,omitempty"` // "auto", "none", "required", 或 ToolChoiceObject
	ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
	StreamOptions       *StreamOptions  `json:"stream_options,omitempty"`
	User                string          `json:"user,omitempty"`
	// 扩展字段
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`
	Seed             *int     `json:"seed,omitempty"`
	N                *int     `json:"n,omitempty"`
}

// ChatMessage 消息结构
type ChatMessage struct {
	Role       string     `json:"role"` // system, user, assistant, tool
	Content    any        `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`  // assistant 消息中的工具调用
	ToolCallID string     `json:"tool_call_id,omitempty"` // tool 消息中的工具调用 ID
}

// ContentPart 多模态内容部分
type ContentPart struct {
	Type     string    `json:"type"` // "text" 或 "image_url"
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL 图片 URL 结构
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // "auto", "low", "high"
}

// ChatTool 工具定义
type ChatTool struct {
	Type     string      `json:"type"` // "function"
	Function FunctionDef `json:"function"`
}

// FunctionDef 函数定义
type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Strict      *bool          `json:"strict,omitempty"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // "function"
	Function FunctionCall `json:"function"`
	Index    *int         `json:"index,omitempty"` // 流式响应中使用
}

// FunctionCall 函数调用详情
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON 字符串
}

// ToolChoiceObject 指定工具选择
type ToolChoiceObject struct {
	Type     string              `json:"type"` // "function"
	Function *ToolChoiceFunction `json:"function,omitempty"`
}

// ToolChoiceFunction 指定函数
type ToolChoiceFunction struct {
	Name string `json:"name"`
}

// ResponseFormat 响应格式
type ResponseFormat struct {
	Type       string          `json:"type"` // "text", "json_object", "json_schema"
	JSONSchema *JSONSchemaSpec `json:"json_schema,omitempty"`
}

// JSONSchemaSpec JSON Schema 规格
type JSONSchemaSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema"`
	Strict      *bool          `json:"strict,omitempty"`
}

// StreamOptions 流式选项
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

// ============================================================
// Response Types
// ============================================================

// ChatCompletionsResponse 非流式响应
type ChatCompletionsResponse struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"` // "chat.completion"
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Choices           []ChatChoice `json:"choices"`
	Usage             *ChatUsage   `json:"usage,omitempty"`
	SystemFingerprint string       `json:"system_fingerprint,omitempty"`
}

// ChatCompletionsChunk 流式响应块
type ChatCompletionsChunk struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"` // "chat.completion.chunk"
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Choices           []ChatChoice `json:"choices"`
	Usage             *ChatUsage   `json:"usage,omitempty"` // 仅在 stream_options.include_usage=true 时
	SystemFingerprint string       `json:"system_fingerprint,omitempty"`
}

// ChatChoice 选择项
type ChatChoice struct {
	Index        int          `json:"index"`
	Message      *ChatMessage `json:"message,omitempty"` // 非流式
	Delta        *ChatDelta   `json:"delta,omitempty"`   // 流式
	FinishReason *string      `json:"finish_reason"`     // "stop", "length", "tool_calls", "content_filter"
	Logprobs     any          `json:"logprobs,omitempty"`
}

// ChatDelta 流式增量
type ChatDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatUsage 用量统计
type ChatUsage struct {
	PromptTokens            int                       `json:"prompt_tokens"`
	CompletionTokens        int                       `json:"completion_tokens"`
	TotalTokens             int                       `json:"total_tokens"`
	PromptTokensDetails     *PromptTokensDetails      `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *CompletionTokensDetails  `json:"completion_tokens_details,omitempty"`
}

// PromptTokensDetails 输入 token 详情
type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
	AudioTokens  int `json:"audio_tokens,omitempty"`
}

// CompletionTokensDetails 输出 token 详情
type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	AudioTokens              int `json:"audio_tokens,omitempty"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

// ============================================================
// Error Types
// ============================================================

// ChatErrorResponse OpenAI 错误响应格式
type ChatErrorResponse struct {
	Error ChatError `json:"error"`
}

// ChatError 错误详情
type ChatError struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param,omitempty"`
	Code    *string `json:"code,omitempty"`
}

// ============================================================
// Helper Methods
// ============================================================

// GetMaxTokens 获取最大输出 token 数
// 优先使用 max_completion_tokens，兼容 max_tokens
func (r *ChatCompletionsRequest) GetMaxTokens() int {
	if r.MaxCompletionTokens != nil && *r.MaxCompletionTokens > 0 {
		return *r.MaxCompletionTokens
	}
	if r.MaxTokens != nil && *r.MaxTokens > 0 {
		return *r.MaxTokens
	}
	return 4096 // 默认值
}

// GetContentAsString 获取消息内容的字符串形式
func (m *ChatMessage) GetContentAsString() string {
	if m.Content == nil {
		return ""
	}
	switch v := m.Content.(type) {
	case string:
		return v
	case []any:
		// 多模态内容，提取文本部分
		var texts []string
		for _, part := range v {
			if partMap, ok := part.(map[string]any); ok {
				if partMap["type"] == "text" {
					if text, ok := partMap["text"].(string); ok {
						texts = append(texts, text)
					}
				}
			}
		}
		return joinStrings(texts, "\n")
	default:
		// 尝试 JSON 序列化
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
		return ""
	}
}

// GetContentParts 获取消息内容的部分列表
func (m *ChatMessage) GetContentParts() []ContentPart {
	if m.Content == nil {
		return nil
	}
	switch v := m.Content.(type) {
	case string:
		return []ContentPart{{Type: "text", Text: v}}
	case []any:
		var parts []ContentPart
		for _, part := range v {
			if partMap, ok := part.(map[string]any); ok {
				cp := ContentPart{}
				if t, ok := partMap["type"].(string); ok {
					cp.Type = t
				}
				if text, ok := partMap["text"].(string); ok {
					cp.Text = text
				}
				if imgURL, ok := partMap["image_url"].(map[string]any); ok {
					cp.ImageURL = &ImageURL{}
					if url, ok := imgURL["url"].(string); ok {
						cp.ImageURL.URL = url
					}
					if detail, ok := imgURL["detail"].(string); ok {
						cp.ImageURL.Detail = detail
					}
				}
				parts = append(parts, cp)
			}
		}
		return parts
	default:
		return nil
	}
}

// HasToolCalls 检查是否有工具调用
func (m *ChatMessage) HasToolCalls() bool {
	return len(m.ToolCalls) > 0
}

// IsToolResult 检查是否是工具结果消息
func (m *ChatMessage) IsToolResult() bool {
	return m.Role == "tool" && m.ToolCallID != ""
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
