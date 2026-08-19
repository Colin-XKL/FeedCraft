package util

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	BrowserProviderBrowserlessRESTful = "browserless-restful"
	BrowserProviderCDP                = "cdp"

	DefaultBrowserRenderTimeout  = 60 * time.Second
	DefaultBrowserMaxConcurrency = 2
	defaultBrowserQueueWait      = 20 * time.Second
)

type BrowserRenderReq struct {
	URL                 string           `json:"url"`
	RejectResourceTypes []string         `json:"rejectResourceTypes,omitempty"`
	WaitForSelector     *WaitForSelector `json:"waitForSelector,omitempty"`
	GotoOptions         *GotoOptions     `json:"gotoOptions,omitempty"`
	WaitFor             int              `json:"waitFor,omitempty"`
}

type WaitForSelector struct {
	Selector  string `json:"selector"`
	TimeoutMs int64  `json:"timeout"`
}

type GotoOptions struct {
	WaitUntil string `json:"waitUntil,omitempty"`
	Timeout   int64  `json:"timeout,omitempty"`
}

type BrowserlessOptions struct {
	Timeout   time.Duration
	WaitTime  time.Duration
	WaitUntil string
}

type BrowserProviderConfig struct {
	Provider string
	Endpoint string
}

// GetBrowserlessContent fetches rendered HTML using the configured browser provider.
func GetBrowserlessContent(websiteUrl string, options BrowserlessOptions) (string, error) {
	envClient := GetEnvClient()
	cfg := ResolveBrowserProviderConfig(envClient)
	if cfg.Endpoint == "" {
		logrus.Errorf("browser provider endpoint not found in env")
		return "", fmt.Errorf("browser provider not configured (FC_BROWSER_ENDPOINT or FC_PUPPETEER_HTTP_ENDPOINT missing)")
	}

	_, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

	if options.Timeout <= 0 {
		options.Timeout = ResolveBrowserRenderTimeout()
	}

	queueWait := defaultBrowserQueueWait
	if options.Timeout > 0 && options.Timeout < queueWait {
		queueWait = options.Timeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), queueWait)
	defer cancel()
	release, err := globalBrowserRenderGate.Acquire(ctx, ResolveBrowserMaxConcurrency())
	if err != nil {
		return "", err
	}
	defer release()

	switch cfg.Provider {
	case BrowserProviderBrowserlessRESTful, "browserless", "":
		return getBrowserlessRESTContent(cfg.Endpoint, websiteUrl, options)
	case BrowserProviderCDP:
		return getCDPContent(cfg.Endpoint, websiteUrl, options)
	default:
		return "", fmt.Errorf("unsupported browser provider %q", cfg.Provider)
	}
}

func ResolveBrowserRenderTimeout() time.Duration {
	raw := strings.TrimSpace(GetEnvClient().GetString("BROWSER_TIMEOUT"))
	if raw == "" {
		return DefaultBrowserRenderTimeout
	}
	if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
		return parsed
	}
	if millis, err := strconv.ParseInt(raw, 10, 64); err == nil && millis > 0 {
		return time.Duration(millis) * time.Millisecond
	}
	logrus.Warnf("Invalid FC_BROWSER_TIMEOUT %q; using default %s", raw, DefaultBrowserRenderTimeout)
	return DefaultBrowserRenderTimeout
}

func ResolveBrowserMaxConcurrency() int {
	limit := GetEnvClient().GetInt("BROWSER_MAX_CONCURRENCY")
	if limit <= 0 {
		return DefaultBrowserMaxConcurrency
	}
	return limit
}

type browserRenderGate struct {
	mu       sync.Mutex
	inflight int
}

func (g *browserRenderGate) Acquire(ctx context.Context, limit int) (func(), error) {
	if g == nil {
		return func() {}, nil
	}
	if limit <= 0 {
		limit = DefaultBrowserMaxConcurrency
	}

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()
	for {
		g.mu.Lock()
		if g.inflight < limit {
			g.inflight++
			g.mu.Unlock()
			return func() {
				g.mu.Lock()
				g.inflight--
				g.mu.Unlock()
			}, nil
		}
		g.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("browser render queue is full: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}

var globalBrowserRenderGate = &browserRenderGate{}

func resetBrowserRenderGateForTest() {
	globalBrowserRenderGate = &browserRenderGate{}
}
