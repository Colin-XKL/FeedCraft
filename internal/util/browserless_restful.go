package util

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

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
