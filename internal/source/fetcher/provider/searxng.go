package provider

import (
	"FeedCraft/internal/config"
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
)

func init() {
	Register("searxng", NewSearXNGProvider)
}

type SearXNGProvider struct {
	Config *config.SearchProviderConfig
}

func NewSearXNGProvider(cfg *config.SearchProviderConfig) SearchProvider {
	return &SearXNGProvider{Config: cfg}
}

func (p *SearXNGProvider) Fetch(ctx context.Context, query string) ([]byte, error) {
	if p.Config.APIUrl == "" {
		return nil, fmt.Errorf("search provider API URL is not configured")
	}

	baseURL := p.Config.APIUrl
	// Remove trailing slash if present
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}

	// Prepare URL parameters
	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")

	// If tool name is provided, treat it as 'engines' (e.g. 'google', 'bing')
	if p.Config.SearchToolName != "" {
		params.Add("engines", p.Config.SearchToolName)
	}

	fullURL := fmt.Sprintf("%s/search?%s", baseURL, params.Encode())

	client := resty.New()
	req := client.R().SetContext(ctx)

	// Add Authorization header if API Key is present (useful for private instances with Basic/Bearer auth)
	if p.Config.APIKey != "" {
		req.SetHeader("Authorization", "Bearer "+p.Config.APIKey)
	}

	resp, err := req.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("search provider returned status %d: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}

func (p *SearXNGProvider) GetDefaultParserConfig() *config.JsonParserConfig {
	return &config.JsonParserConfig{
		ItemsIterator: "results",
		Title:         "title",
		Link:          "url",
		Description:   "content",
		Date:          "publishedDate",
	}
}
