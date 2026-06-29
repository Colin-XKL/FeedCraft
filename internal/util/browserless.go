package util

import (
	"FeedCraft/internal/config"
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
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
	Timeout           time.Duration
	WaitTime          time.Duration
	WaitUntil         string
	NavigationActions []config.BrowserNavigationAction
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
	if err := ValidateBrowserNavigationActions(options.NavigationActions); err != nil {
		return "", err
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
