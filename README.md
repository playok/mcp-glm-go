[English](./README.md) | [한국어](./README_ko.md) | [日本語](./README_jp.md) | [中文](./README_cn.md)

# mcp-glm-go

An MCP (Model Context Protocol) server that provides Z.AI GLM model capabilities as tools.
Use GLM chat, thinking, web search, and image generation in Google Antigravity IDE via `@mcp:glm:`.

## Tools

| Tool | Description |
|------|-------------|
| `glm_chat` | General chat, code generation, text creation |
| `glm_chat_with_thinking` | Thinking mode - returns reasoning process and final answer |
| `glm_web_search` | Web search enabled - answers with up-to-date info and sources |
| `glm_image_gen` | Image generation with CogView-4 / GLM-Image |

## Installation

### Requirements
- Go 1.21+
- Z.AI API key (get one at [z.ai](https://z.ai))

### Build

```bash
git clone https://github.com/playok/mcp-glm-go.git
cd mcp-glm-go
go build -o mcp-glm-go .
```

## Antigravity Configuration

Add the MCP server to your Antigravity settings. You can register both GLM-4.7 and GLM-5 to choose by use case.

```json
{
  "mcpServers": {
    "glm47": {
      "command": "/absolute/path/to/mcp-glm-go",
      "args": ["--coding"],
      "env": {
        "GLM_API_KEY": "your-api-key"
      }
    },
    "glm5": {
      "command": "/absolute/path/to/mcp-glm-go",
      "args": ["--coding", "--model", "glm-5"],
      "env": {
        "GLM_API_KEY": "your-api-key"
      }
    }
  }
}
```

- `@mcp:glm47:` — GLM-4.7 (Coding Plan default, fast and lightweight)
- `@mcp:glm5:` — GLM-5 (745B flagship, ideal for complex reasoning)

### Single Model Setup

GLM-4.7 only:

```json
{
  "mcpServers": {
    "glm": {
      "command": "/absolute/path/to/mcp-glm-go",
      "args": ["--coding"],
      "env": {
        "GLM_API_KEY": "your-api-key"
      }
    }
  }
}
```

GLM-5 only:

```json
{
  "mcpServers": {
    "glm": {
      "command": "/absolute/path/to/mcp-glm-go",
      "args": ["--coding", "--model", "glm-5"],
      "env": {
        "GLM_API_KEY": "your-api-key"
      }
    }
  }
}
```

### CLI Options

| Option | Description | Default |
|--------|-------------|---------|
| `--api-key` | API key (also via `GLM_API_KEY` env, flag takes priority) | - |
| `--coding` | Use Coding Plan endpoint (`api/coding/paas/v4`) | false |
| `--model` | Default chat model | `glm-4.7` |

## Usage

Call tools in Antigravity chat using `@mcp:<server>:<tool>` format.

With both models registered (`glm47`, `glm5`):

```
# Chat with GLM-4.7
@mcp:glm47:glm_chat Implement quicksort in Python

# Chat with GLM-5
@mcp:glm5:glm_chat Implement quicksort in Python

# GLM-5 Thinking mode (complex reasoning)
@mcp:glm5:glm_chat_with_thinking Analyze the time complexity of this algorithm

# GLM-4.7 Web search
@mcp:glm47:glm_web_search What are the latest Go 1.24 changes?

# Image generation
@mcp:glm47:glm_image_gen A futuristic city skyline, cyberpunk style
```

With single model (`glm`):

```
@mcp:glm:glm_chat Implement quicksort in Python
@mcp:glm:glm_chat_with_thinking Analyze the time complexity of this algorithm
@mcp:glm:glm_web_search What are the latest Go 1.24 changes?
@mcp:glm:glm_image_gen A futuristic city skyline, cyberpunk style
```

### Parameters

Each tool supports additional parameters:

- **glm_chat**: `prompt`, `model`, `system_msg`, `temperature`, `max_tokens`
- **glm_chat_with_thinking**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_web_search**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_image_gen**: `prompt`, `model`(cogview-4-250304/glm-image), `size`, `quality`(hd/standard)

## Supported Models

| Model | Description | Note |
|-------|-------------|------|
| `glm-5` | Latest flagship (745B MoE) | Standard API, paid |
| `glm-4.7` | Coding-optimized | Included in Coding Plan |
| `glm-4.7-flash` | Lightweight free model | Standard API, free |
| `cogview-4-250304` | Image generation | $0.01/image |
| `glm-image` | Image generation (high-res) | $0.015/image |

## License

MIT
