[English](./README.md) | [한국어](./README_ko.md) | [日本語](./README_jp.md) | [中文](./README_cn.md)

# mcp-glm-go

Z.AIのGLMモデルをMCP（Model Context Protocol）ツールとして提供するサーバーです。
Google Antigravity IDEで `@mcp:glm:` 形式でGLMモデルのチャット、思考（thinking）、Web検索、画像生成機能を使用できます。

## 提供ツール

| ツール | 説明 |
|--------|------|
| `glm_chat` | 一般的な会話、コード生成、テキスト作成 |
| `glm_chat_with_thinking` | 思考モード - 推論過程と最終回答を返却 |
| `glm_web_search` | Web検索有効化 - 最新情報とソースを含む回答 |
| `glm_image_gen` | CogView-4 / GLM-Imageによる画像生成 |

## インストール

### 要件
- Go 1.21以上
- Z.AI APIキー（[z.ai](https://z.ai)で取得）

### ビルド

```bash
git clone https://github.com/playok/mcp-glm-go.git
cd mcp-glm-go
go build -o mcp-glm-go .
```

## APIキー設定（オプション）

Antigravity設定ファイルにAPIキーを記載する代わりに、システム環境変数として設定できます。

**macOS / Linux (bash/zsh)**

```bash
# ~/.bashrc、~/.zshrc、または ~/.bash_profile に追加
export GLM_API_KEY="your-api-key"

# 即座に反映
source ~/.zshrc  # または ~/.bashrc
```

**Windows (PowerShell)**

```powershell
# 現在のセッションのみ
$env:GLM_API_KEY = "your-api-key"

# 永続設定（ユーザーレベル）
[System.Environment]::SetEnvironmentVariable("GLM_API_KEY", "your-api-key", "User")
```

**Windows (CMD)**

```cmd
# 現在のセッションのみ
set GLM_API_KEY=your-api-key

# 永続設定（ユーザーレベル）
setx GLM_API_KEY "your-api-key"
```

`GLM_API_KEY` 環境変数を設定すると、Antigravity設定で `env` ブロックを省略できます：

```json
{
  "mcpServers": {
    "glm": {
      "command": "/absolute/path/to/mcp-glm-go",
      "args": ["--coding"]
    }
  }
}
```

## Antigravity設定

Antigravityの設定ファイルにMCPサーバーを追加します。GLM-4.7とGLM-5を同時に登録して、用途に応じて使い分けることができます。

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

- `@mcp:glm47:` — GLM-4.7モデル（Coding Planデフォルト、高速で軽量）
- `@mcp:glm5:` — GLM-5モデル（745Bフラッグシップ、複雑な推論に最適）

### 単一モデルのみ使用する場合

GLM-4.7のみ：

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

GLM-5のみ：

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

### CLIオプション

| オプション | 説明 | デフォルト |
|------------|------|------------|
| `--api-key` | APIキー（環境変数 `GLM_API_KEY` でも設定可能、フラグ優先） | - |
| `--coding` | Coding Planエンドポイントを使用（`api/coding/paas/v4`） | false |
| `--model` | デフォルトチャットモデルを指定 | `glm-4.7` |

> **セキュリティ注意:** `--api-key` フラグを使用すると、`ps aux` などのプロセス一覧でAPIキーが表示される可能性があります。`GLM_API_KEY` 環境変数の使用を推奨します。

## 使い方

Antigravityチャットで `@mcp:<サーバー名>:<ツール名>` 形式で呼び出します。

両モデル登録の場合（`glm47`、`glm5`）：

```
# GLM-4.7でチャット
@mcp:glm47:glm_chat Pythonでクイックソートを実装して

# GLM-5でチャット
@mcp:glm5:glm_chat Pythonでクイックソートを実装して

# GLM-5 Thinkingモード（複雑な推論）
@mcp:glm5:glm_chat_with_thinking このアルゴリズムの時間計算量を分析して

# GLM-4.7 Web検索
@mcp:glm47:glm_web_search 2026年最新のGo 1.24の変更点を教えて

# 画像生成
@mcp:glm47:glm_image_gen 未来都市の風景、サイバーパンクスタイル
```

単一モデル登録の場合（`glm`）：

```
@mcp:glm:glm_chat Pythonでクイックソートを実装して
@mcp:glm:glm_chat_with_thinking このアルゴリズムの時間計算量を分析して
@mcp:glm:glm_web_search 2026年最新のGo 1.24の変更点を教えて
@mcp:glm:glm_image_gen 未来都市の風景、サイバーパンクスタイル
```

### パラメータ

各ツールは追加パラメータをサポートしています：

- **glm_chat**: `prompt`, `model`, `system_msg`, `temperature`, `max_tokens`
- **glm_chat_with_thinking**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_web_search**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_image_gen**: `prompt`, `model`(cogview-4-250304/glm-image), `size`, `quality`(hd/standard)

## 対応モデル

| モデル | 説明 | 備考 |
|--------|------|------|
| `glm-5` | 最新フラッグシップ (745B MoE) | 標準API、有料 |
| `glm-4.7` | コーディング最適化モデル | Coding Plan含む |
| `glm-4.7-flash` | 軽量無料モデル | 標準API、無料 |
| `cogview-4-250304` | 画像生成 | $0.01/枚 |
| `glm-image` | 画像生成（高解像度） | $0.015/枚 |

## ライセンス

MIT
