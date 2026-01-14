package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/util"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HttpFetcher is a simple fetcher based on http.Get.
type HttpFetcher struct {
	Config *config.HttpFetcherConfig
}

func (f *HttpFetcher) Fetch(ctx context.Context) ([]byte, error) {
	if f.Config == nil || f.Config.URL == "" {
		return nil, fmt.Errorf("http fetcher is not configured with a URL")
	}

	if f.Config.UseBrowserless {
		content, err := util.GetBrowserlessContent(f.Config.URL, util.BrowserlessOptions{
			Timeout: 30 * time.Second,
		}) // TODO: Make timeout configurable?
		if err != nil {
			return nil, fmt.Errorf("browserless fetch failed: %w", err)
		}
		return []byte(content), nil
	}

	method := f.Config.Method
	if method == "" {
		method = "GET"
	}

	var bodyReader io.Reader
	if f.Config.Body != "" {
		bodyReader = strings.NewReader(f.Config.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, f.Config.URL, bodyReader)
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
	defer func() { _ = resp.Body.Close() }()

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
