package provider

import (
	"FeedCraft/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func init() {
	Register("litellm", NewLiteLLMProvider)
}

type LiteLLMProvider struct {
	Config *config.SearchProviderConfig
}

func NewLiteLLMProvider(cfg *config.SearchProviderConfig) SearchProvider {
	return &LiteLLMProvider{Config: cfg}
}

func (p *LiteLLMProvider) Fetch(ctx context.Context, query string) ([]byte, error) {
	if p.Config.APIUrl == "" {
		return nil, fmt.Errorf("search provider API URL is not configured")
	}

	// Construct URL
	url := p.Config.APIUrl
	// LiteLLM Search usually maps to a specific model/tool name if using OpenAI format,
	// but strictly speaking, the "Search" endpoint might be different.
	// Based on previous code: if SearchToolName is present, append it.
	if p.Config.SearchToolName != "" {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		url += p.Config.SearchToolName
	}

	// Prepare Request Body
	reqBody := map[string]interface{}{
		"query": query,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.Config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.Config.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search provider returned status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (p *LiteLLMProvider) GetDefaultParserConfig() *config.JsonParserConfig {
	return &config.JsonParserConfig{
		ItemsIterator: "data",
		Title:         "title",
		Link:          "url",
		Description:   "content",
		// Date might vary, usually not standard in simple search results
	}
}
