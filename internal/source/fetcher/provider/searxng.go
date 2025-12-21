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
	Client *resty.Client
}

func NewSearXNGProvider(cfg *config.SearchProviderConfig) SearchProvider {
	return &SearXNGProvider{
		Config: cfg,
		Client: resty.New(),
	}
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

	// If engines is provided, treat it as 'engines' (e.g. 'google', 'bing')
	if p.Config.SearXNG.Engines != "" {
		params.Add("engines", p.Config.SearXNG.Engines)
	}

	fullURL := fmt.Sprintf("%s/search?%s", baseURL, params.Encode())

	req := p.Client.R().SetContext(ctx)

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
