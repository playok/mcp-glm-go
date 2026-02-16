package main

// GLM Chat API types

type GLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GLMThinking struct {
	Type string `json:"type"` // "enabled" or "disabled"
}

type GLMWebSearch struct {
	Enable        bool   `json:"enable"`
	SearchEngine  string `json:"search_engine,omitempty"`
	SearchResult  string `json:"search_result,omitempty"`
}

type GLMChatRequest struct {
	Model      string        `json:"model"`
	Messages   []GLMMessage  `json:"messages"`
	Stream     bool          `json:"stream"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens  *int          `json:"max_tokens,omitempty"`
	Thinking   *GLMThinking  `json:"thinking,omitempty"`
	WebSearch  *GLMWebSearch `json:"web_search,omitempty"`
}

type GLMChoice struct {
	Index        int        `json:"index"`
	FinishReason string     `json:"finish_reason"`
	Message      GLMResMsg  `json:"message"`
}

type GLMResMsg struct {
	Role             string             `json:"role"`
	Content          string             `json:"content"`
	ReasoningContent string             `json:"reasoning_content,omitempty"`
	ToolCalls        []GLMToolCall      `json:"tool_calls,omitempty"`
}

type GLMToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function GLMFunctionCall `json:"function"`
}

type GLMFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type GLMUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GLMWebSearchResult struct {
	Refer []GLMSearchRefer `json:"refer,omitempty"`
}

type GLMSearchRefer struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Content string `json:"content"`
}

type GLMChatResponse struct {
	ID        string          `json:"id"`
	Created   int64           `json:"created"`
	Model     string          `json:"model"`
	Choices   []GLMChoice     `json:"choices"`
	Usage     GLMUsage        `json:"usage"`
	WebSearch *GLMWebSearchResult `json:"web_search,omitempty"`
}

// GLM Image Generation API types

type GLMImageRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Quality string `json:"quality,omitempty"`
	Size    string `json:"size,omitempty"`
}

type GLMImageData struct {
	URL string `json:"url"`
}

type GLMImageResponse struct {
	Created int64          `json:"created"`
	Data    []GLMImageData `json:"data"`
}

// GLM Error types

type GLMErrorResponse struct {
	Error GLMErrorDetail `json:"error"`
}

type GLMErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
