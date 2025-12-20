package provider

import (
	"FeedCraft/internal/config"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header if API Key is present (useful for private instances with Basic/Bearer auth)
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

func (p *SearXNGProvider) GetDefaultParserConfig() *config.JsonParserConfig {
	return &config.JsonParserConfig{
		ItemsIterator: "results",
		Title:         "title",
		Link:          "url",
		Description:   "content",
		Date:          "publishedDate",
	}
}
