package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input types

type ChatInput struct {
	Prompt      string   `json:"prompt" jsonschema:"user message to send"`
	Model       string   `json:"model,omitempty" jsonschema:"model ID override (default: glm-4.7)"`
	SystemMsg   string   `json:"system_msg,omitempty" jsonschema:"optional system message"`
	Temperature *float64 `json:"temperature,omitempty" jsonschema:"sampling temperature between 0.0 and 1.0"`
	MaxTokens   *int     `json:"max_tokens,omitempty" jsonschema:"maximum number of output tokens"`
}

type ThinkingInput struct {
	Prompt    string `json:"prompt" jsonschema:"user message to send"`
	Model     string `json:"model,omitempty" jsonschema:"model ID override (default: glm-4.7)"`
	SystemMsg string `json:"system_msg,omitempty" jsonschema:"optional system message"`
	MaxTokens *int   `json:"max_tokens,omitempty" jsonschema:"maximum number of output tokens"`
}

type WebSearchInput struct {
	Prompt    string `json:"prompt" jsonschema:"question to search and answer"`
	Model     string `json:"model,omitempty" jsonschema:"model ID override (default: glm-4.7)"`
	SystemMsg string `json:"system_msg,omitempty" jsonschema:"optional system message"`
	MaxTokens *int   `json:"max_tokens,omitempty" jsonschema:"maximum number of output tokens"`
}

type ImageGenInput struct {
	Prompt  string `json:"prompt" jsonschema:"description of the image to generate"`
	Model   string `json:"model,omitempty" jsonschema:"model to use: cogview-4-250304 or glm-image (default: cogview-4-250304)"`
	Size    string `json:"size,omitempty" jsonschema:"image size e.g. 1024x1024"`
	Quality string `json:"quality,omitempty" jsonschema:"quality: hd or standard"`
}

const defaultChatModel = "glm-4.7"

func registerTools(server *mcp.Server, client *GLMClient, defaultModel string) {
	if defaultModel == "" {
		defaultModel = defaultChatModel
	}
	registerChatTool(server, client, defaultModel)
	registerThinkingTool(server, client, defaultModel)
	registerWebSearchTool(server, client, defaultModel)
	registerImageGenTool(server, client)
}

func resolveModel(input string, defaultModel string) string {
	if input != "" {
		return input
	}
	return defaultModel
}

func validateChatParams(temperature *float64, maxTokens *int) error {
	if temperature != nil && (*temperature < 0.0 || *temperature > 1.0) {
		return fmt.Errorf("temperature must be between 0.0 and 1.0")
	}
	if maxTokens != nil && *maxTokens < 1 {
		return fmt.Errorf("max_tokens must be at least 1")
	}
	return nil
}

func registerChatTool(server *mcp.Server, client *GLMClient, defaultModel string) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "glm_chat",
		Description: "Chat with GLM model. Use for general questions, code generation, and text creation. | GLM 모델과 대화합니다. | GLMモデルと会話します。 | 与GLM模型对话，用于问答、代码生成和文本创作。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ChatInput) (*mcp.CallToolResult, any, error) {
		if input.Prompt == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "prompt is required. | prompt는 필수입니다. | promptは必須です。 | prompt为必填项。"}},
				IsError: true,
			}, nil, nil
		}
		if err := validateChatParams(input.Temperature, input.MaxTokens); err != nil {
			return errorResult(err), nil, nil
		}

		messages := buildMessages(input.SystemMsg, input.Prompt)
		chatReq := &GLMChatRequest{
			Model:       resolveModel(input.Model, defaultModel),
			Messages:    messages,
			Stream:      false,
			Temperature: input.Temperature,
			MaxTokens:   input.MaxTokens,
		}

		resp, err := client.ChatCompletion(ctx, chatReq)
		if err != nil {
			return errorResult(err), nil, nil
		}

		content := extractContent(resp)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: content}},
		}, nil, nil
	})
}

func registerThinkingTool(server *mcp.Server, client *GLMClient, defaultModel string) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "glm_chat_with_thinking",
		Description: "Chat with thinking mode enabled. Returns both reasoning process and final answer. | 사고 모드로 추론 과정과 최종 답변을 반환합니다. | 思考モードで推論過程と最終回答を返します。 | 启用思考模式，返回推理过程和最终答案。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ThinkingInput) (*mcp.CallToolResult, any, error) {
		if input.Prompt == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "prompt is required. | prompt는 필수입니다. | promptは必須です。 | prompt为必填项。"}},
				IsError: true,
			}, nil, nil
		}
		if err := validateChatParams(nil, input.MaxTokens); err != nil {
			return errorResult(err), nil, nil
		}

		messages := buildMessages(input.SystemMsg, input.Prompt)
		chatReq := &GLMChatRequest{
			Model:    resolveModel(input.Model, defaultModel),
			Messages: messages,
			Stream:   false,
			Thinking: &GLMThinking{Type: "enabled"},
			MaxTokens: input.MaxTokens,
		}

		resp, err := client.ChatCompletion(ctx, chatReq)
		if err != nil {
			return errorResult(err), nil, nil
		}

		var parts []string
		if len(resp.Choices) > 0 {
			msg := resp.Choices[0].Message
			if msg.ReasoningContent != "" {
				parts = append(parts, fmt.Sprintf("<thinking>\n%s\n</thinking>", msg.ReasoningContent))
			}
			if msg.Content != "" {
				parts = append(parts, msg.Content)
			}
		}

		text := strings.Join(parts, "\n\n")
		if text == "" {
			text = "(empty response)"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: text}},
		}, nil, nil
	})
}

func registerWebSearchTool(server *mcp.Server, client *GLMClient, defaultModel string) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "glm_web_search",
		Description: "Chat with web search enabled. Returns up-to-date answers with source references. | 웹 검색으로 최신 정보와 출처를 포함한 답변을 생성합니다. | Web検索で最新情報とソースを含む回答を生成します。 | 启用网络搜索，返回包含来源的最新信息。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input WebSearchInput) (*mcp.CallToolResult, any, error) {
		if input.Prompt == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "prompt is required. | prompt는 필수입니다. | promptは必須です。 | prompt为必填项。"}},
				IsError: true,
			}, nil, nil
		}
		if err := validateChatParams(nil, input.MaxTokens); err != nil {
			return errorResult(err), nil, nil
		}

		messages := buildMessages(input.SystemMsg, input.Prompt)
		chatReq := &GLMChatRequest{
			Model:    resolveModel(input.Model, defaultModel),
			Messages: messages,
			Stream:   false,
			WebSearch: &GLMWebSearch{
				Enable:       true,
				SearchEngine: "search_pro_jina",
			},
			MaxTokens: input.MaxTokens,
		}

		resp, err := client.ChatCompletion(ctx, chatReq)
		if err != nil {
			return errorResult(err), nil, nil
		}

		var parts []string
		content := extractContent(resp)
		parts = append(parts, content)

		if resp.WebSearch != nil && len(resp.WebSearch.Refer) > 0 {
			parts = append(parts, "\n---\n**Sources:**")
			for i, ref := range resp.WebSearch.Refer {
				parts = append(parts, fmt.Sprintf("%d. [%s](%s)", i+1, ref.Title, ref.Link))
			}
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(parts, "\n")}},
		}, nil, nil
	})
}

func registerImageGenTool(server *mcp.Server, client *GLMClient) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "glm_image_gen",
		Description: "Generate images with CogView-4 or GLM-Image. Returns inline image and URL. | 이미지를 생성하고 인라인 이미지와 URL을 반환합니다. | 画像を生成し、インライン画像とURLを返します。 | 生成图像并返回内联图像和URL。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ImageGenInput) (*mcp.CallToolResult, any, error) {
		if input.Prompt == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "prompt is required. | prompt는 필수입니다. | promptは必須です。 | prompt为必填项。"}},
				IsError: true,
			}, nil, nil
		}

		model := input.Model
		if model == "" {
			model = "cogview-4-250304"
		}

		imgReq := &GLMImageRequest{
			Model:   model,
			Prompt:  input.Prompt,
			Size:    input.Size,
			Quality: input.Quality,
		}

		resp, err := client.ImageGeneration(ctx, imgReq)
		if err != nil {
			return errorResult(err), nil, nil
		}

		if len(resp.Data) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "No image generated. | 이미지 생성 결과가 없습니다. | 画像が生成されませんでした。 | 未生成图像。"}},
				IsError: true,
			}, nil, nil
		}

		imageURL := resp.Data[0].URL
		var contents []mcp.Content

		// 이미지 다운로드 후 인라인 반환 시도
		imgData, mimeType, err := client.DownloadImage(ctx, imageURL)
		if err != nil {
			log.Printf("image download failed, returning URL only: %v", err)
		} else {
			contents = append(contents, &mcp.ImageContent{
				Data:     imgData,
				MIMEType: mimeType,
			})
		}

		// URL도 항상 포함 (30일간 유효)
		contents = append(contents, &mcp.TextContent{
			Text: fmt.Sprintf("Image URL (valid for 30 days): %s", imageURL),
		})

		return &mcp.CallToolResult{Content: contents}, nil, nil
	})
}

// Helper functions

func buildMessages(systemMsg, prompt string) []GLMMessage {
	var messages []GLMMessage
	if systemMsg != "" {
		messages = append(messages, GLMMessage{Role: "system", Content: systemMsg})
	}
	messages = append(messages, GLMMessage{Role: "user", Content: prompt})
	return messages
}

func extractContent(resp *GLMChatResponse) string {
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		return resp.Choices[0].Message.Content
	}
	return "(empty response)"
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
		IsError: true,
	}
}
