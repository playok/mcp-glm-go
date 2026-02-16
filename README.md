# GLM MCP Server

Z.AI의 GLM 모델을 MCP(Model Context Protocol) 도구로 제공하는 서버입니다.
Google Antigravity IDE에서 `@mcp:glm:` 형태로 GLM 모델의 채팅, 사고(thinking), 웹 검색, 이미지 생성 기능을 사용할 수 있습니다.

## 제공 도구

| 도구 | 설명 |
|------|------|
| `glm_chat` | 일반 대화, 코드 작성, 텍스트 생성 |
| `glm_chat_with_thinking` | 사고(thinking) 모드 - 추론 과정과 최종 답변 반환 |
| `glm_web_search` | 웹 검색 활성화 - 최신 정보 + 출처 포함 답변 |
| `glm_image_gen` | CogView-4 / GLM-Image 이미지 생성 |

## 설치

### 요구사항
- Go 1.21 이상
- Z.AI API 키 ([z.ai](https://z.ai)에서 발급)

### 빌드

```bash
git clone https://github.com/playok/mcp-glm-go.git
cd mcp-glm-go
go build -o mcp-glm-go .
```

## Antigravity 구성

Antigravity 설정 파일에 MCP 서버를 추가합니다. glm-4.7과 glm-5를 동시에 등록하면 용도에 따라 골라 사용할 수 있습니다.

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

- `@mcp:glm47:` — GLM-4.7 모델 (Coding Plan 기본, 빠르고 가벼움)
- `@mcp:glm5:` — GLM-5 모델 (745B 플래그십, 복잡한 추론에 적합)

### 단일 모델만 사용할 경우

GLM-4.7만 사용:

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

GLM-5만 사용:

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

### CLI 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--api-key` | API 키 (환경변수 `GLM_API_KEY`로도 설정 가능, 플래그 우선) | - |
| `--coding` | Coding Plan 엔드포인트 사용 (`api/coding/paas/v4`) | false |
| `--model` | 기본 채팅 모델 지정 | `glm-4.7` |

## 사용법

Antigravity 채팅에서 `@mcp:<서버이름>:<도구이름>` 형태로 호출합니다.

두 모델을 모두 등록한 경우 (`glm47`, `glm5`):

```
# GLM-4.7로 대화
@mcp:glm47:glm_chat 파이썬으로 퀵소트 구현해줘

# GLM-5로 대화
@mcp:glm5:glm_chat 파이썬으로 퀵소트 구현해줘

# GLM-5 Thinking 모드 (복잡한 추론)
@mcp:glm5:glm_chat_with_thinking 이 알고리즘의 시간 복잡도를 분석해줘

# GLM-4.7 웹 검색
@mcp:glm47:glm_web_search 2026년 최신 Go 1.24 변경사항 알려줘

# 이미지 생성
@mcp:glm47:glm_image_gen 미래 도시 풍경, 사이버펑크 스타일
```

단일 모델 등록인 경우 (`glm`):

```
@mcp:glm:glm_chat 파이썬으로 퀵소트 구현해줘
@mcp:glm:glm_chat_with_thinking 이 알고리즘의 시간 복잡도를 분석해줘
@mcp:glm:glm_web_search 2026년 최신 Go 1.24 변경사항 알려줘
@mcp:glm:glm_image_gen 미래 도시 풍경, 사이버펑크 스타일
```

### 파라미터

각 도구는 추가 파라미터를 지원합니다:

- **glm_chat**: `prompt`, `model`, `system_msg`, `temperature`, `max_tokens`
- **glm_chat_with_thinking**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_web_search**: `prompt`, `model`, `system_msg`, `max_tokens`
- **glm_image_gen**: `prompt`, `model`(cogview-4-250304/glm-image), `size`, `quality`(hd/standard)

## 지원 모델

| 모델 | 설명 | 비고 |
|------|------|------|
| `glm-5` | 최신 플래그십 (745B MoE) | 일반 API, 유료 |
| `glm-4.7` | Coding 최적화 모델 | Coding Plan 포함 |
| `glm-4.7-flash` | 경량 무료 모델 | 일반 API, 무료 |
| `cogview-4-250304` | 이미지 생성 | $0.01/장 |
| `glm-image` | 이미지 생성 (고해상도) | $0.015/장 |

## 라이선스

MIT
