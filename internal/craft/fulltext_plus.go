package craft

import (
	"FeedCraft/internal/util"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
)

type browserRenderReq struct {
	URL                 string           `json:"url"`
	RejectResourceTypes []string         `json:"rejectResourceTypes,omitempty"`
	WaitForSelector     *WaitForSelector `json:"waitForSelector,omitempty"`
	//BestAttempt         bool             `json:"bestAttempt"`
}
type WaitForSelector struct {
	Selector  string `json:"selector"`
	TimeoutMs int64  `json:"timeout"`
}

func getRenderedHTML2(websiteUrl string, timeout time.Duration) (string, error) {
	envClient := util.GetEnvClient()
	browserURI := envClient.GetString("PUPPETEER_HTTP_ENDPOINT")
	if browserURI == "" {
		log.Fatalf("puppeteer websocket endpoint PUPPETEER_HTTP_ENDPOINT not found in env")
	}
	parseUrl, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

	// Use resty client with timeout and retry
	client := resty.New().
		SetBaseURL(browserURI).
		SetTimeout(timeout).
		SetRetryCount(util.MaxRetries).
		SetRetryWaitTime(util.RetryWaitTime).
		SetRetryMaxWaitTime(util.RetryMaxWaitTime)

	headers := map[string]string{
		"Cache-Control": "no-cache",
		"Content-Type":  "application/json",
	}
	reqBody := browserRenderReq{
		URL:                 websiteUrl,
		RejectResourceTypes: []string{"image"},
		//BestAttempt:         true,
		//WaitForSelector: &WaitForSelector{
		//	Selector:  "body",
		//	TimeoutMs: 30000,
		//},
	}
	response, err := client.R().SetHeaders(headers).SetBody(reqBody).Post("/content")
	if err != nil {
		return "", err
	}
	//fmt.Println(response.String())

	article, err := readability.FromReader(strings.NewReader(response.String()), parseUrl)
	if err != nil {
		return "", err
	}
	return article.Content, err
}

func GetFulltextPlusCraftOptions() []CraftOption {
	transFunc := func(item *feeds.Item) (string, error) {
		link := item.Link.Href
		return getRenderedHTML2(link, util.ExternalRequestTimeout)
	}

	cachedTransFunc := GetCommonCachedTransformer(cacheKeyForArticleLink, transFunc, "extract fulltext plus")

	relativeLinkFixOptions := GetRelativeLinkFixCraftOptions()

	var craftOptions []CraftOption
	craftOptions = append(craftOptions, relativeLinkFixOptions...)
	craftOptions = append(craftOptions, OptionTransformFeedItem(GetArticleContentProcessor(cachedTransFunc)))
	return craftOptions
}
