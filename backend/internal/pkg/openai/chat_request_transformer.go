package openai

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// ============================================================
// OpenAI Chat -> Claude Messages 请求转换
// ============================================================

// ClaudeMessagesRequest Claude Messages API 请求格式
// 用于转换后直接序列化发送
type ClaudeMessagesRequest struct {
	Model       string              `json:"model"`
	Messages    []ClaudeMessage     `json:"messages"`
	MaxTokens   int                 `json:"max_tokens"`
	System      any                 `json:"system,omitempty"` // string 或 []SystemBlock
	Stream      bool                `json:"stream,omitempty"`
	Temperature *float64            `json:"temperature,omitempty"`
	TopP        *float64            `json:"top_p,omitempty"`
	Tools       []ClaudeTool        `json:"tools,omitempty"`
	Metadata    *ClaudeMetadata     `json:"metadata,omitempty"`
}

// ClaudeMessage Claude 消息
type ClaudeMessage struct {
	Role    string `json:"role"` // user, assistant
	Content any    `json:"content"`
}

// SystemBlock system prompt 数组元素
type SystemBlock struct {
	Type         string            `json:"type"`
	Text         string            `json:"text"`
	CacheControl *CacheControlSpec `json:"cache_control,omitempty"`
}

// CacheControlSpec 缓存控制
type CacheControlSpec struct {
	Type string `json:"type"` // "ephemeral"
}

// ClaudeTool Claude 工具定义
type ClaudeTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema,omitempty"`
}

// ClaudeMetadata Claude 请求元数据
type ClaudeMetadata struct {
	UserID string `json:"user_id,omitempty"`
}

// ClaudeContentBlock Claude 内容块
type ClaudeContentBlock struct {
	Type      string `json:"type"`                 // text, tool_use, tool_result, image
	Text      string `json:"text,omitempty"`       // type=text
	ID        string `json:"id,omitempty"`         // type=tool_use
	Name      string `json:"name,omitempty"`       // type=tool_use
	Input     any    `json:"input,omitempty"`      // type=tool_use
	ToolUseID string `json:"tool_use_id,omitempty"` // type=tool_result
	Content   any    `json:"content,omitempty"`    // type=tool_result (string 或 []ContentBlock)
	IsError   bool   `json:"is_error,omitempty"`   // type=tool_result
	Source    *ImageSource `json:"source,omitempty"` // type=image
}

// ImageSource 图片来源
type ImageSource struct {
	Type      string `json:"type"`       // "base64" 或 "url"
	MediaType string `json:"media_type,omitempty"`
	Data      string `json:"data,omitempty"`
	URL       string `json:"url,omitempty"`
}

// TransformChatToClaude 将 OpenAI Chat Completions 请求转换为 Claude Messages 格式
func TransformChatToClaude(req *ChatCompletionsRequest, enableCache bool) (*ClaudeMessagesRequest, error) {
	claudeReq := &ClaudeMessagesRequest{
		Model:       req.Model,
		MaxTokens:   req.GetMaxTokens(),
		Stream:      req.Stream,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	// 1. 转换消息
	system, messages, err := transformMessages(req.Messages, enableCache)
	if err != nil {
		return nil, fmt.Errorf("transform messages: %w", err)
	}
	claudeReq.System = system
	claudeReq.Messages = messages

	// 2. 转换工具定义
	if len(req.Tools) > 0 {
		claudeReq.Tools = transformTools(req.Tools)
	}

	// 3. 设置 metadata（用于粘性会话）
	if req.User != "" {
		claudeReq.Metadata = &ClaudeMetadata{
			UserID: req.User,
		}
	}

	return claudeReq, nil
}

// transformMessages 转换消息列表
// 返回 system prompt 和 Claude 消息列表
func transformMessages(messages []ChatMessage, enableCache bool) (any, []ClaudeMessage, error) {
	var systemParts []SystemBlock
	var claudeMessages []ClaudeMessage

	// 用于收集连续的 tool 消息
	var pendingToolResults []ClaudeContentBlock

	for i, msg := range messages {
		switch msg.Role {
		case "system":
			// system 消息转换为 system prompt
			text := msg.GetContentAsString()
			block := SystemBlock{
				Type: "text",
				Text: text,
			}
			// 第一个 system block 添加缓存控制
			if enableCache && len(systemParts) == 0 {
				block.CacheControl = &CacheControlSpec{Type: "ephemeral"}
			}
			systemParts = append(systemParts, block)

		case "user":
			// 先 flush 之前的 tool results
			if len(pendingToolResults) > 0 {
				claudeMessages = append(claudeMessages, ClaudeMessage{
					Role:    "user",
					Content: pendingToolResults,
				})
				pendingToolResults = nil
			}

			// user 消息
			content, err := transformUserContent(msg)
			if err != nil {
				return nil, nil, fmt.Errorf("transform user message %d: %w", i, err)
			}
			claudeMessages = append(claudeMessages, ClaudeMessage{
				Role:    "user",
				Content: content,
			})

		case "assistant":
			// 先 flush 之前的 tool results
			if len(pendingToolResults) > 0 {
				claudeMessages = append(claudeMessages, ClaudeMessage{
					Role:    "user",
					Content: pendingToolResults,
				})
				pendingToolResults = nil
			}

			// assistant 消息
			content := transformAssistantContent(msg)
			claudeMessages = append(claudeMessages, ClaudeMessage{
				Role:    "assistant",
				Content: content,
			})

		case "tool":
			// tool 消息转换为 tool_result，收集起来
			toolResult := ClaudeContentBlock{
				Type:      "tool_result",
				ToolUseID: msg.ToolCallID,
				Content:   msg.GetContentAsString(),
			}
			pendingToolResults = append(pendingToolResults, toolResult)
		}
	}

	// flush 剩余的 tool results
	if len(pendingToolResults) > 0 {
		claudeMessages = append(claudeMessages, ClaudeMessage{
			Role:    "user",
			Content: pendingToolResults,
		})
	}

	// 处理 system prompt
	var system any
	if len(systemParts) == 1 && systemParts[0].CacheControl == nil {
		// 单个 system 且没有缓存控制，使用字符串格式
		system = systemParts[0].Text
	} else if len(systemParts) > 0 {
		// 多个 system 或有缓存控制，使用数组格式
		system = systemParts
	}

	return system, claudeMessages, nil
}

// transformUserContent 转换 user 消息内容
func transformUserContent(msg ChatMessage) (any, error) {
	parts := msg.GetContentParts()
	if len(parts) == 0 {
		return msg.GetContentAsString(), nil
	}

	// 检查是否只有文本
	allText := true
	for _, p := range parts {
		if p.Type != "text" {
			allText = false
			break
		}
	}

	if allText && len(parts) == 1 {
		return parts[0].Text, nil
	}

	// 转换为 Claude content blocks
	var blocks []ClaudeContentBlock
	for _, p := range parts {
		switch p.Type {
		case "text":
			blocks = append(blocks, ClaudeContentBlock{
				Type: "text",
				Text: p.Text,
			})
		case "image_url":
			if p.ImageURL != nil {
				block := ClaudeContentBlock{Type: "image"}
				// 检查是否是 base64 data URL
				if isBase64DataURL(p.ImageURL.URL) {
					mediaType, data := parseBase64DataURL(p.ImageURL.URL)
					block.Source = &ImageSource{
						Type:      "base64",
						MediaType: mediaType,
						Data:      data,
					}
				} else {
					// URL 格式
					block.Source = &ImageSource{
						Type: "url",
						URL:  p.ImageURL.URL,
					}
				}
				blocks = append(blocks, block)
			}
		}
	}

	return blocks, nil
}

// transformAssistantContent 转换 assistant 消息内容
func transformAssistantContent(msg ChatMessage) any {
	// 如果有 tool_calls，需要转换为 content blocks
	if msg.HasToolCalls() {
		var blocks []ClaudeContentBlock

		// 先添加文本内容（如果有）
		text := msg.GetContentAsString()
		if text != "" {
			blocks = append(blocks, ClaudeContentBlock{
				Type: "text",
				Text: text,
			})
		}

		// 添加 tool_use blocks
		for _, tc := range msg.ToolCalls {
			var input any
			if tc.Function.Arguments != "" {
				_ = json.Unmarshal([]byte(tc.Function.Arguments), &input)
			}
			blocks = append(blocks, ClaudeContentBlock{
				Type:  "tool_use",
				ID:    tc.ID,
				Name:  tc.Function.Name,
				Input: input,
			})
		}

		return blocks
	}

	// 纯文本
	return msg.GetContentAsString()
}

// transformTools 转换工具定义
func transformTools(tools []ChatTool) []ClaudeTool {
	claudeTools := make([]ClaudeTool, 0, len(tools))
	for _, t := range tools {
		if t.Type == "function" {
			claudeTools = append(claudeTools, ClaudeTool{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				InputSchema: t.Function.Parameters,
			})
		}
	}
	return claudeTools
}

// ============================================================
// Session Hash 生成
// ============================================================

var sessionIDRegex = regexp.MustCompile(`session_([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)

// ExtractSessionID 从 user 字段提取 session ID
func ExtractSessionID(user string) string {
	if match := sessionIDRegex.FindStringSubmatch(user); len(match) > 1 {
		return match[1]
	}
	return ""
}

// ============================================================
// Helper Functions
// ============================================================

var base64DataURLRegex = regexp.MustCompile(`^data:([^;]+);base64,(.+)$`)

func isBase64DataURL(url string) bool {
	return base64DataURLRegex.MatchString(url)
}

func parseBase64DataURL(url string) (mediaType, data string) {
	matches := base64DataURLRegex.FindStringSubmatch(url)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return "", ""
}

// SerializeClaudeRequest 序列化 Claude 请求为 JSON
func SerializeClaudeRequest(req *ClaudeMessagesRequest) ([]byte, error) {
	return json.Marshal(req)
}
