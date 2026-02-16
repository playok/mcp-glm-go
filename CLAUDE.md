# mcp-glm-go

MCP server for Z.AI GLM models, written in Go.

## Project Structure

```
├── main.go          # Entrypoint: CLI flags (--api-key, --coding, --model), server setup, stdio transport
├── types.go         # GLM API request/response structs (chat, image, error)
├── glmclient.go     # HTTP client: ChatCompletion, ImageGeneration, DownloadImage
├── tools.go         # 4 MCP tools: glm_chat, glm_chat_with_thinking, glm_web_search, glm_image_gen
├── README.md        # English (default)
├── README_ko.md     # Korean
├── README_jp.md     # Japanese
├── README_cn.md     # Chinese (Simplified)
└── .github/workflows/release.yml  # Multi-platform release (mac/linux/windows)
```

## Key Design Decisions

- **Default model**: `glm-4.7` (changeable via `--model` flag)
- **Coding Plan endpoint**: `--coding` flag switches to `api/coding/paas/v4`
- **Stdio transport**: MCP communication via stdin/stdout JSON-RPC (persistent process)
- **Closure pattern**: Tool handlers capture `*GLMClient` via closures
- **Image handling**: Downloads generated URL → base64 `ImageContent` + URL fallback
- **Multilingual**: Tool descriptions in EN/KO/JA/ZH separated by `|`

## API Endpoints

| Endpoint | URL |
|----------|-----|
| Standard | `https://api.z.ai/api/paas/v4` |
| Coding Plan | `https://api.z.ai/api/coding/paas/v4` |

## Build & Release

```bash
# Local build
go build -o mcp-glm-go .

# Release: push a tag to trigger GitHub Actions
git tag v0.x.0 && git push origin v0.x.0
```

## Testing

```bash
# MCP handshake test (stdin must stay open with sleep)
{
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'
  echo '{"jsonrpc":"2.0","method":"notifications/initialized"}'
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
  sleep 5
} | ./mcp-glm-go --api-key YOUR_KEY --coding 2>/dev/null
```
