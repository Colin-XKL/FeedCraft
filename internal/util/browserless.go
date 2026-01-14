package util

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
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

// GetBrowserlessContent fetches the rendered HTML content of a URL using the browserless service.
// It relies on the PUPPETEER_HTTP_ENDPOINT environment variable.
func GetBrowserlessContent(websiteUrl string, options BrowserlessOptions) (string, error) {
	envClient := GetEnvClient()
	browserURI := envClient.GetString("PUPPETEER_HTTP_ENDPOINT")
	if browserURI == "" {
		// Log warning instead of fatal, as this might be called in contexts where we want to handle the error
		logrus.Errorf("puppeteer websocket endpoint PUPPETEER_HTTP_ENDPOINT not found in env")
		return "", fmt.Errorf("browserless service not configured (PUPPETEER_HTTP_ENDPOINT missing)")
	}
	// Since we are moving to a utility, returning an error is better.
	// But if the env is missing, it's a configuration error.
	// I'll stick to error return.

	_, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

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
		return "", fmt.Errorf("browserless service returned status %d: %s", response.StatusCode(), response.String())
	}

	return response.String(), nil
}
