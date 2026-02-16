package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	log.SetOutput(os.Stderr)

	apiKey := flag.String("api-key", "", "GLM API key (also via GLM_API_KEY env)")
	coding := flag.Bool("coding", false, "Use GLM Coding Plan endpoint (api/coding/paas/v4)")
	model := flag.String("model", "glm-4.7", "Default chat model (e.g. glm-4.7, glm-5, glm-4.7-flash)")
	flag.Parse()

	key := *apiKey
	if key == "" {
		key = os.Getenv("GLM_API_KEY")
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "API key required: use --api-key flag or GLM_API_KEY env variable")
		os.Exit(1)
	}

	baseURL := ""
	if *coding {
		baseURL = codingBaseURL
	}
	client := NewGLMClient(key, baseURL)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "glm-mcp-server",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Instructions: "GLM MCP Server: chat, thinking, web search, image generation tools for Z.AI GLM models.",
		},
	)

	registerTools(server, client, *model)

	log.Println("GLM MCP Server starting on stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
