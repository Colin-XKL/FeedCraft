package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	BrowserProviderBrowserlessRESTful = "browserless-restful"
	BrowserProviderCDP                = "cdp"
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
		options.Timeout = 30 * time.Second
	}

	switch cfg.Provider {
	case BrowserProviderBrowserlessRESTful, "browserless", "":
		return getBrowserlessRESTContent(cfg.Endpoint, websiteUrl, options)
	case BrowserProviderCDP:
		return getCDPContent(cfg.Endpoint, websiteUrl, options)
	default:
		return "", fmt.Errorf("unsupported browser provider %q", cfg.Provider)
	}
}

func IsSupportedBrowserProvider(provider string) bool {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case BrowserProviderBrowserlessRESTful, "browserless", BrowserProviderCDP:
		return true
	default:
		return false
	}
}

func ResolveBrowserProviderConfig(env *viper.Viper) BrowserProviderConfig {
	provider := strings.ToLower(strings.TrimSpace(env.GetString("BROWSER_PROVIDER")))
	endpoint := strings.TrimSpace(env.GetString("BROWSER_ENDPOINT"))
	if endpoint == "" {
		endpoint = strings.TrimSpace(env.GetString("PUPPETEER_HTTP_ENDPOINT"))
	}
	if provider == "" {
		provider = BrowserProviderBrowserlessRESTful
	}
	return BrowserProviderConfig{
		Provider: provider,
		Endpoint: endpoint,
	}
}

func resolveBrowserProviderConfig(env *viper.Viper) BrowserProviderConfig {
	return ResolveBrowserProviderConfig(env)
}

func getBrowserlessRESTContent(browserURI string, websiteUrl string, options BrowserlessOptions) (string, error) {
	client := resty.New().SetBaseURL(browserURI)
	client.SetTimeout(options.Timeout)

	headers := map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  "application/json",
	}
	reqBody := BrowserRenderReq{
		URL:                 websiteUrl,
		RejectResourceTypes: []string{"image"},
		WaitFor:             int(options.WaitTime.Milliseconds()),
	}

	if options.WaitUntil != "" {
		reqBody.GotoOptions = &GotoOptions{
			WaitUntil: options.WaitUntil,
		}
	}

	response, err := client.R().SetHeaders(headers).SetBody(reqBody).Post("/content")
	if err != nil {
		return "", err
	}

	if response.StatusCode() != http.StatusOK {
		respStr := response.String()
		logrus.Errorf("browserless service returned status %d. URL: %s, response body: %s", response.StatusCode(), websiteUrl, respStr)

		truncLen := 200
		if len(respStr) > truncLen {
			respStr = respStr[:truncLen] + "..."
		}
		return "", fmt.Errorf("browserless service returned status %d: %s", response.StatusCode(), respStr)
	}

	return response.String(), nil
}

func getCDPContent(endpoint string, websiteUrl string, options BrowserlessOptions) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
	defer cancel()

	wsURL, err := getCDPWebSocketURL(ctx, endpoint)
	if err != nil {
		return "", err
	}

	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(ctx, wsURL)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	idleMonitor := newNetworkIdleMonitor()
	chromedp.ListenTarget(browserCtx, idleMonitor.handleEvent)

	var html string
	actions := []chromedp.Action{
		network.Enable(),
		network.SetBlockedURLs([]string{
			"*.png", "*.jpg", "*.jpeg", "*.gif", "*.webp", "*.svg", "*.ico",
		}),
		navigateAndWaitForPageEventAction(websiteUrl, options.WaitUntil),
	}
	if strings.EqualFold(options.WaitUntil, "networkidle0") {
		actions = append(actions, idleMonitor.waitAction(0, 500*time.Millisecond))
	} else if strings.EqualFold(options.WaitUntil, "networkidle2") {
		actions = append(actions, idleMonitor.waitAction(2, 500*time.Millisecond))
	}
	if options.WaitTime > 0 {
		actions = append(actions, chromedp.Sleep(options.WaitTime))
	}
	actions = append(actions, chromedp.Evaluate(`document.documentElement.outerHTML`, &html))

	if err := chromedp.Run(browserCtx, actions...); err != nil {
		return "", fmt.Errorf("browser cdp render failed: %w", err)
	}
	return html, nil
}

func navigateAndWaitForPageEventAction(websiteURL string, waitUntil string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		waitForDomContent := strings.EqualFold(waitUntil, "domcontentloaded") ||
			strings.EqualFold(waitUntil, "networkidle0") ||
			strings.EqualFold(waitUntil, "networkidle2")
		done := make(chan struct{})
		var once sync.Once
		markDone := func() {
			once.Do(func() {
				close(done)
			})
		}

		chromedp.ListenTarget(ctx, func(ev any) {
			switch ev.(type) {
			case *page.EventDomContentEventFired:
				if waitForDomContent {
					markDone()
				}
			case *page.EventLoadEventFired:
				if !waitForDomContent {
					markDone()
				}
			}
		})

		if err := page.Enable().Do(ctx); err != nil {
			return err
		}
		_, _, errorText, _, err := page.Navigate(websiteURL).Do(ctx)
		if err != nil {
			return err
		}
		if errorText != "" {
			return fmt.Errorf("navigation failed: %s", errorText)
		}

		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type cdpVersionResponse struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func getCDPWebSocketURL(ctx context.Context, endpoint string) (string, error) {
	versionURL, err := BuildEndpointURL(endpoint, "/json/version")
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, versionURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("browser cdp version request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 200))
		return "", fmt.Errorf("browser cdp service returned status %d: %s", resp.StatusCode, string(body))
	}

	var version cdpVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return "", fmt.Errorf("failed to decode browser cdp version response: %w", err)
	}
	if version.WebSocketDebuggerURL == "" {
		return "", fmt.Errorf("browser cdp version response missing webSocketDebuggerUrl")
	}
	return version.WebSocketDebuggerURL, nil
}

func BuildEndpointURL(endpoint string, path string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	basePath := strings.TrimRight(u.Path, "/")
	nextPath := "/" + strings.TrimLeft(path, "/")
	if basePath == "" {
		u.Path = nextPath
	} else {
		u.Path = basePath + nextPath
	}
	return u.String(), nil
}

func buildEndpointURL(endpoint string, path string) (string, error) {
	return BuildEndpointURL(endpoint, path)
}

type networkIdleMonitor struct {
	mu       sync.Mutex
	inflight map[network.RequestID]struct{}
}

func newNetworkIdleMonitor() *networkIdleMonitor {
	return &networkIdleMonitor{inflight: make(map[network.RequestID]struct{})}
}

func (m *networkIdleMonitor) handleEvent(ev any) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch e := ev.(type) {
	case *network.EventRequestWillBeSent:
		m.inflight[e.RequestID] = struct{}{}
	case *network.EventLoadingFinished:
		delete(m.inflight, e.RequestID)
	case *network.EventLoadingFailed:
		delete(m.inflight, e.RequestID)
	}
}

func (m *networkIdleMonitor) waitAction(maxInflight int, idleFor time.Duration) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var idleSince time.Time
		for {
			m.mu.Lock()
			inflight := len(m.inflight)
			m.mu.Unlock()

			if inflight <= maxInflight {
				if idleSince.IsZero() {
					idleSince = time.Now()
				}
				if time.Since(idleSince) >= idleFor {
					return nil
				}
			} else {
				idleSince = time.Time{}
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}
		}
	}
}
