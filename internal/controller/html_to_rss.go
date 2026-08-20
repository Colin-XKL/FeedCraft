package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/craft"
	fetcherpkg "FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/source/parser"
	"FeedCraft/internal/util"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type FetchReq struct {
	URL            string `json:"url" binding:"required"`
	UseBrowserless bool   `json:"use_browserless"`
}

type ParseReq struct {
	HTML            string `json:"html"`
	URL             string `json:"url"` // Fallback if HTML is empty, fetch internally
	ItemSelector    string `json:"item_selector"`
	TitleSelector   string `json:"title_selector"`
	LinkSelector    string `json:"link_selector"`
	DateSelector    string `json:"date_selector"`
	ContentSelector string `json:"content_selector"`
}

type ParsedItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

type WebMonitorPreviewReq struct {
	HTML             string                        `json:"html"`
	URL              string                        `json:"url" binding:"required"`
	UseBrowserless   bool                          `json:"use_browserless"`
	WebMonitorParser config.WebMonitorParserConfig `json:"web_monitor_parser"`
}

func validateURL(rawUrl string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return fmt.Errorf("access to private IP %s is forbidden", ip.String())
		}
	}
	return nil
}

// fetchHTML extracts common fetching logic with browser emulation and error handling
func fetchHTML(targetURL string, useBrowserless bool) (string, error) {
	if useBrowserless {
		return util.GetBrowserlessContent(targetURL, util.BrowserlessOptions{
			Timeout: util.ResolveBrowserRenderTimeout(),
		})
	}

	client := resty.New()
	client.SetTimeout(craft.DefaultExtractFulltextTimeout)
	req := client.R()
	for key, value := range fetcherpkg.HTMLDefaultHeaders() {
		req.SetHeader(key, value)
	}
	resp, err := req.Get(targetURL)

	if err != nil {
		return "", fmt.Errorf("fetch failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("upstream returned status %d. The site might be blocking bots. Try enabling 'Enhance Mode'", resp.StatusCode())
	}

	content := resp.String()
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("upstream returned 200 OK but the content is empty. Try enabling 'Enhance Mode'")
	}

	return content, nil
}

func HtmlFetch(c *gin.Context) {
	var req FetchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	if err := validateURL(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	htmlContent, err := fetchHTML(req.URL, req.UseBrowserless)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[string]{
		StatusCode: 0,
		Data:       htmlContent,
	})
}

func HtmlParse(c *gin.Context) {
	var req ParseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	var htmlContent string
	var err error

	if req.HTML != "" {
		htmlContent = req.HTML
	} else if req.URL != "" {
		if err := validateURL(req.URL); err != nil {
			c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
			return
		}

		htmlContent, err = fetchHTML(req.URL, false)
		if err != nil {
			c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: "Either html or url is required"})
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Parse HTML failed: " + err.Error()})
		return
	}

	var items []ParsedItem
	if req.ItemSelector == "" {
		c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{StatusCode: 0, Data: items})
		return
	}

	doc.Find(req.ItemSelector).Each(func(i int, s *goquery.Selection) {
		item := ParsedItem{}

		getSelection := func(selector string) *goquery.Selection {
			if selector == "" || selector == "." {
				return s
			}
			return s.Find(selector)
		}

		if req.TitleSelector != "" {
			item.Title = strings.TrimSpace(getSelection(req.TitleSelector).Text())
		}
		if req.LinkSelector != "" {
			sel := getSelection(req.LinkSelector)
			item.Link = util.ExtractLinkFromSelection(sel)

			if req.URL != "" && item.Link != "" {
				if absURL, err := util.BuildAbsoluteURL(req.URL, item.Link); err == nil {
					item.Link = absURL
				}
			}

			if item.Link != "" {
				if u, err := url.Parse(item.Link); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
					item.Link = ""
				}
			}
		}
		if req.DateSelector != "" {
			sel := getSelection(req.DateSelector)
			item.Date = strings.TrimSpace(sel.Text())
			if item.Date == "" {
				val, exists := sel.Attr("datetime")
				if exists {
					item.Date = val
				}
			}
		}
		if req.ContentSelector != "" {
			sel := getSelection(req.ContentSelector)
			htmlStr, err := sel.Html()
			if err != nil {
				logrus.Infof("Warning: Failed to extract content using selector '%s' for item %d in feed %s: %v",
					req.ContentSelector, i, req.URL, err)
			}
			if err == nil && htmlStr != "" {
				item.Content = htmlStr
			}
		}

		items = append(items, item)
	})

	c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{
		StatusCode: 0,
		Data:       items,
	})
}

func WebMonitorPreview(c *gin.Context) {
	var req WebMonitorPreviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	if err := validateURL(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	htmlContent := req.HTML
	if strings.TrimSpace(htmlContent) == "" {
		var err error
		htmlContent, err = fetchHTML(req.URL, req.UseBrowserless)
		if err != nil {
			c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
			return
		}
	}

	preview, err := parser.PreviewWebMonitor([]byte(htmlContent), &req.WebMonitorParser, req.URL)
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[*parser.WebMonitorPreviewResult]{
		StatusCode: 0,
		Data:       preview,
	})
}
