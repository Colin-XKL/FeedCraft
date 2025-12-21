package provider

import (
	"FeedCraft/internal/config"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func init() {
	Register("litellm", NewLiteLLMProvider)
}

type LiteLLMProvider struct {
	Config *config.SearchProviderConfig
	Client *resty.Client
}

func NewLiteLLMProvider(cfg *config.SearchProviderConfig) SearchProvider {
	return &LiteLLMProvider{
		Config: cfg,
		Client: resty.New().SetTimeout(10 * time.Second),
	}
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
	if p.Config.LiteLLM.SearchToolName != "" {
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}
		url += p.Config.LiteLLM.SearchToolName
	}

	// Lazy initialization for safety
	if p.Client == nil {
		p.Client = resty.New().SetTimeout(10 * time.Second)
	}

	// Prepare Request Body
	reqBody := map[string]interface{}{
		"query": query,
	}

	req := p.Client.R().
		SetContext(ctx).
		SetBody(reqBody)

	if p.Config.APIKey != "" {
		req.SetHeader("Authorization", "Bearer "+p.Config.APIKey)
	}

	// Resty automatically sets Content-Type to application/json when SetBody is used with a map/struct

	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("search provider returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
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