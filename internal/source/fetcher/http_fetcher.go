package fetcher

import (
	"FeedCraft/internal/config"
	"context"
	"fmt"
	"io"
	"net/http"
)

// HttpFetcher is a simple fetcher based on http.Get.
type HttpFetcher struct {
	Config *config.HttpFetcherConfig
}

func (f *HttpFetcher) Fetch(ctx context.Context) ([]byte, error) {
	if f.Config == nil || f.Config.URL == "" {
		return nil, fmt.Errorf("http fetcher is not configured with a URL")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", f.Config.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// It is good practice to set a user-agent.
	req.Header.Set("User-Agent", "FeedCraft/2.0")

	// Set custom headers from config
	for key, value := range f.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status not ok: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (f *HttpFetcher) BaseURL() string {
	if f.Config == nil {
		return ""
	}
	return f.Config.URL
}
