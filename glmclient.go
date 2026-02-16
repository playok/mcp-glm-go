package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxAPIResponseSize = 10 << 20 // 10MB for chat/image API responses
	maxImageDownload   = 50 << 20 // 50MB for image file downloads
)

const (
	defaultBaseURL = "https://api.z.ai/api/paas/v4"
	codingBaseURL  = "https://api.z.ai/api/coding/paas/v4"
)

type GLMClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewGLMClient(apiKey, baseURL string) *GLMClient {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &GLMClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *GLMClient) ChatCompletion(ctx context.Context, req *GLMChatRequest) (*GLMChatResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxAPIResponseSize))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	var chatResp GLMChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &chatResp, nil
}

func (c *GLMClient) ImageGeneration(ctx context.Context, req *GLMImageRequest) (*GLMImageResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/images/generations", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxAPIResponseSize))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	var imgResp GLMImageResponse
	if err := json.Unmarshal(respBody, &imgResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &imgResp, nil
}

func (c *GLMClient) DownloadImage(ctx context.Context, imageURL string) ([]byte, string, error) {
	if err := validateImageURL(imageURL); err != nil {
		return nil, "", fmt.Errorf("invalid image URL: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("create download request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download image returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxImageDownload))
	if err != nil {
		return nil, "", fmt.Errorf("read image data: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/png"
	}

	return data, mimeType, nil
}

// validateImageURL checks that the URL is safe to fetch (SSRF prevention).
func validateImageURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("malformed URL")
	}

	if u.Scheme != "https" {
		return fmt.Errorf("only https URLs are allowed")
	}

	host := u.Hostname()

	// Block localhost and loopback
	if host == "localhost" || strings.HasPrefix(host, "127.") || host == "::1" {
		return fmt.Errorf("loopback addresses are not allowed")
	}

	// Block private/internal IP ranges
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("private/internal addresses are not allowed")
		}
	}

	return nil
}

// parseAPIError extracts a safe error message from API error responses.
func parseAPIError(statusCode int, body []byte) error {
	var errResp GLMErrorResponse
	if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
		return fmt.Errorf("GLM API error (%s): %s", errResp.Error.Code, errResp.Error.Message)
	}
	return fmt.Errorf("GLM API returned status %d", statusCode)
}
