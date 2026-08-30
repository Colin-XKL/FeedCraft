package util

import (
	"FeedCraft/internal/config"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type BrowserFunctionReq struct {
	Code    string                 `json:"code"`
	Context BrowserFunctionContext `json:"context"`
}

type BrowserFunctionContext struct {
	URL               string                           `json:"url"`
	WaitUntil         string                           `json:"wait_until,omitempty"`
	WaitTimeMs        int64                            `json:"wait_time_ms,omitempty"`
	TimeoutMs         int64                            `json:"timeout_ms,omitempty"`
	NavigationActions []config.BrowserNavigationAction `json:"navigation_actions,omitempty"`
}

func getBrowserlessRESTContent(browserURI string, websiteUrl string, options BrowserlessOptions) (string, error) {
	client := resty.New().SetBaseURL(browserURI)
	client.SetTimeout(options.Timeout)

	headers := map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  "application/json",
	}
	if len(options.NavigationActions) > 0 {
		reqBody := BrowserFunctionReq{
			Code: browserNavigationFunctionCode,
			Context: BrowserFunctionContext{
				URL:               websiteUrl,
				WaitUntil:         options.WaitUntil,
				WaitTimeMs:        options.WaitTime.Milliseconds(),
				TimeoutMs:         options.Timeout.Milliseconds(),
				NavigationActions: options.NavigationActions,
			},
		}
		return postBrowserlessRequest(client, headers, reqBody, "/function", websiteUrl)
	}

	reqBody := BrowserRenderReq{
		URL:                 websiteUrl,
		RejectResourceTypes: []string{"image"},
		WaitFor:             int(options.WaitTime.Milliseconds()),
		GotoOptions: &GotoOptions{
			Timeout: options.Timeout.Milliseconds(),
		},
	}

	if options.WaitUntil != "" {
		reqBody.GotoOptions.WaitUntil = options.WaitUntil
	}

	return postBrowserlessRequest(client, headers, reqBody, "/content", websiteUrl)
}

func postBrowserlessRequest(client *resty.Client, headers map[string]string, reqBody any, path string, websiteUrl string) (string, error) {
	response, err := client.R().SetHeaders(headers).SetBody(reqBody).Post(path)
	if err != nil {
		return "", err
	}

	if response.StatusCode() != http.StatusOK {
		respStr := response.String()
		logrus.Errorf("browserless service returned status %d. Path: %s, URL: %s, response body: %s", response.StatusCode(), path, websiteUrl, respStr)

		truncLen := 200
		if len(respStr) > truncLen {
			respStr = respStr[:truncLen] + "..."
		}
		return "", fmt.Errorf("browserless service returned status %d: %s", response.StatusCode(), respStr)
	}

	return response.String(), nil
}

const browserNavigationFunctionCode = `module.exports = async ({ page, context }) => {
  const waitUntil = context.wait_until || "load";
  const timeout = context.timeout_ms || 30000;
  await page.setRequestInterception(true);
  page.on("request", (request) => {
    if (request.resourceType() === "image") {
      request.abort();
      return;
    }
    request.continue();
  });
  await page.goto(context.url, { waitUntil, timeout });
  for (const action of context.navigation_actions || []) {
    if (action.type === "click") {
      await page.waitForSelector(action.selector, { timeout: action.timeout_ms || timeout, visible: true });
      await page.click(action.selector);
    } else if (action.type === "wait_for_selector") {
      await page.waitForSelector(action.selector, { timeout: action.timeout_ms || timeout, visible: true });
    } else if (action.type === "wait") {
      await new Promise((resolve) => setTimeout(resolve, action.duration_ms));
    } else {
      throw new Error("unsupported navigation action: " + action.type);
    }
  }
  if (context.wait_time_ms > 0) {
    await new Promise((resolve) => setTimeout(resolve, context.wait_time_ms));
  }
  return { data: await page.content(), type: "text/html" };
};`
