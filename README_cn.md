[English](./README.md) | [한국어](./README_ko.md) | [日本語](./README_jp.md) | [中文](./README_cn.md)

# mcp-glm-go

将Z.AI的GLM模型作为MCP（Model Context Protocol）工具提供的服务器。
在Google Antigravity IDE中通过 `@mcp:glm:` 格式使用GLM模型的聊天、思考（thinking）、网络搜索和图像生成功能。

## 提供工具

| 工具 | 说明 |
|------|------|
| `glm_chat` | 通用对话、代码生成、文本创作 |
| `glm_chat_with_thinking` | 思考模式 - 返回推理过程和最终答案 |
| `glm_web_search` | 启用网络搜索 - 包含最新信息和来源的回答 |
| `glm_image_gen` | 使用CogView-4 / GLM-Image生成图像 |

## 安装

### 要求
- Go 1.21及以上
- Z.AI API密钥（在 [z.ai](https://z.ai) 获取）

### 构建

```bash
git clone https://github.com/playok/mcp-glm-go.git
cd mcp-glm-go
go build -o mcp-glm-go .
```

## Antigravity配置

在Antigravity设置文件中添加MCP服务器。同时注册GLM-4.7和GLM-5，可根据用途选择使用。

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

- `@mcp:glm47:` — GLM-4.7模型（Coding Plan默认，快速轻量）
- `@mcp:glm5:` — GLM-5模型（745B旗舰，适合复杂推理）

### 仅使用单个模型

仅使用GLM-4.7：

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

仅使用GLM-5：

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

### CLI选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `--api-key` | API密钥（也可通过环境变量 `GLM_API_KEY` 设置，标志优先） | - |
| `--coding` | 使用Coding Plan端点（`api/coding/paas/v4`） | false |
| `--model` | 指定默认聊天模型 | `glm-4.7` |

## 使用方法

在Antigravity聊天中使用 `@mcp:<服务器名>:<工具名>` 格式调用。

同时注册两个模型的情况（`glm47`、`glm5`）：

```
# 使用GLM-4.7聊天
@mcp:glm47:glm_chat 用Python实现快速排序

# 使用GLM-5聊天
@mcp:glm5:glm_chat 用Python实现快速排序

# GLM-5 Thinking模式（复杂推理）
@mcp:glm5:glm_chat_with_thinking 分析这个算法的时间复杂度

# GLM-4.7 网络搜索
@mcp:glm47:glm_web_search 2026年Go 1.24有哪些最新变化

# 图像生成
@mcp:glm47:glm_image_gen 未来城市景观，赛博朋克风格
```

仅注册单个模型的情况（`glm`）：

```
@mcp:glm:glm_chat 用Python实现快速排序
@mcp:glm:glm_chat_with_thinking 分析这个算法的时间复杂度
@mcp:glm:glm_web_search 2026年Go 1.24有哪些最新变化
@mcp:glm:glm_image_gen 未来城市景观，赛博朋克风格
```

### 参数

每个工具支持额外参数：

- **glm_chat**: `prompt`, `model`, `system_msg`, `temperature`, `max_tokens`
- **glm_chat_with_thinking**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_web_search**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_image_gen**: `prompt`, `model`(cogview-4-250304/glm-image), `size`, `quality`(hd/standard)

## 支持模型

| 模型 | 说明 | 备注 |
|------|------|------|
| `glm-5` | 最新旗舰 (745B MoE) | 标准API，付费 |
| `glm-4.7` | 编码优化模型 | Coding Plan包含 |
| `glm-4.7-flash` | 轻量免费模型 | 标准API，免费 |
| `cogview-4-250304` | 图像生成 | $0.01/张 |
| `glm-image` | 图像生成（高分辨率） | $0.015/张 |

## 许可证

MIT
