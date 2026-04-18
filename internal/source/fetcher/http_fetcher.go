package fetcher

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/util"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	retry "github.com/avast/retry-go/v4"
)

type requestProfile struct {
	defaultHeaders map[string]string
	retryAttempts  uint
}

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

	profile := resolveRequestProfile(f.Config)
	var body []byte
	err := retry.Do(
		func() error {
			result, err := f.doRequest(ctx, profile)
			if err != nil {
				return err
			}
			body = result
			return nil
		},
		retry.Attempts(profile.retryAttempts),
		retry.Delay(300*time.Millisecond),
		retry.DelayType(retry.FixedDelay),
		retry.RetryIf(isRetryableFetchError),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (f *HttpFetcher) doRequest(ctx context.Context, profile requestProfile) ([]byte, error) {
	method := f.Config.Method
	if method == "" {
		method = http.MethodGet
	}

	var bodyReader io.Reader
	if f.Config.Body != "" {
		bodyReader = strings.NewReader(f.Config.Body)
	}

	req, err := http.NewRequestWithContext(ctx, method, f.Config.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range profile.defaultHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range f.Config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &fetchError{err: fmt.Errorf("http get failed: %w", err), retryable: true}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, &fetchError{
			err:       fmt.Errorf("http status not ok: %s", resp.Status),
			retryable: isRetryableStatus(resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &fetchError{err: fmt.Errorf("failed to read response body: %w", err), retryable: true}
	}

	return body, nil
}

func (f *HttpFetcher) BaseURL() string {
	if f.Config == nil {
		return ""
	}
	return f.Config.URL
}

func resolveRequestProfile(cfg *config.HttpFetcherConfig) requestProfile {
	if cfg != nil && cfg.Purpose == config.HttpFetcherPurposeHTML {
		return requestProfile{
			defaultHeaders: HTMLDefaultHeaders(),
			retryAttempts:  3,
		}
	}

	return requestProfile{
		defaultHeaders: map[string]string{
			"User-Agent": util.DefaultFeedUserAgent(),
		},
		retryAttempts: 1,
	}
}

func HTMLDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":                util.DefaultHTMLUserAgent(),
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "en-US,en;q=0.9",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
	}
}

type fetchError struct {
	err       error
	retryable bool
}

func (e *fetchError) Error() string {
	return e.err.Error()
}

func (e *fetchError) Unwrap() error {
	return e.err
}

func isRetryableFetchError(err error) bool {
	var fetchErr *fetchError
	if !errors.As(err, &fetchErr) {
		return false
	}
	return fetchErr.retryable
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}
