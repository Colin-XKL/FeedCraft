package recipe

import (
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"strings"
	"time"
)

//func getRenderedHTML(websiteUrl string) (string, error) {
//	envClient := util.GetEnvClient()
//	browserURI := envClient.GetString("PUPPETEER_WS_ENDPOINT")
//	if browserURI == "" {
//		log.Fatalf("puppeteer websocket endpoint not found in env")
//	}
//
//	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), browserURI)
//	defer cancel()
//
//	// create context
//	ctx, cancel := chromedp.NewContext(allocatorContext)
//	defer cancel()
//
//	var res string
//
//	err := chromedp.Run(ctx,
//		chromedp.Navigate(websiteUrl),
//		//chromedp.WaitVisible("body"),
//		chromedp.Sleep(1*time.Second),
//		chromedp.ActionFunc(func(ctx context.Context) error {
//			fmt.Println("reading dom")
//			node, err := dom.GetDocument().Do(ctx)
//			if err != nil {
//				return err
//			}
//			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
//			return err
//		}),
//	)
//
//	if err != nil {
//		logrus.Errorf("get rendered web page failed: %v", err)
//	}
//	return res, nil
//
//}

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
		log.Fatalf("puppeteer websocket endpoint not found in env")
	}
	parseUrl, err := url.Parse(websiteUrl)
	if err != nil {
		logrus.Errorf("parse url failed: %v", err)
		return "", err
	}

	client := resty.New().SetBaseURL(browserURI)
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

	return article.Content, err
	return response.String(), nil
}

func ExtractFulltextPlusForFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(500, "empty feed url")
		return
	}
	fp := gofeed.NewParser()
	parsedFeed, _ := fp.ParseURL(feedUrl)

	ret := TransformFeed(parsedFeed, GetFulltextExtractor(getRenderedHTML2))

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
