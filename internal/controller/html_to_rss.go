package controller

import (
	"FeedCraft/internal/craft"
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
			Timeout: craft.DefaultExtractFulltextTimeout,
		})
	}

	// Try standard HTTP request (simulating a browser user agent)
	client := resty.New()
	client.SetTimeout(craft.DefaultExtractFulltextTimeout)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1").
		Get(targetURL)

	if err != nil {
		return "", fmt.Errorf("fetch failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("upstream returned status %d. The site might be blocking bots. Try enabling 'Enhance Mode'", resp.StatusCode())
	}

	content := resp.String()
	// Check if body is empty even with 200 OK
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
		// Use StatusCode: -1 to indicate logic/upstream error rather than system error
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

		// Fallback fetch if HTML not provided. Default to standard fetch (no browserless) as ParseReq doesn't support it yet.
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
	// If no item selector, return empty
	if req.ItemSelector == "" {
		c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{StatusCode: 0, Data: items})
		return
	}

	doc.Find(req.ItemSelector).Each(func(i int, s *goquery.Selection) {
		item := ParsedItem{}

		// Helper to extract selection based on selector
		// If selector is "." or empty (though frontend sends . now), use current 's'
		// Otherwise find descendant
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

			// Try to resolve relative URL to absolute URL
			if req.URL != "" && item.Link != "" {
				if absURL, err := util.BuildAbsoluteURL(req.URL, item.Link); err == nil {
					item.Link = absURL
				}
			}

			// Final validation: Ensure it is a valid HTTP/HTTPS URL
			if item.Link != "" {
				if u, err := url.Parse(item.Link); err != nil || (u.Scheme != "http" && u.Scheme != "https") {
					item.Link = ""
				}
			}
		}
		if req.DateSelector != "" {
			sel := getSelection(req.DateSelector)
			item.Date = strings.TrimSpace(sel.Text())
			// Check datetime attr if text is empty
			if item.Date == "" {
				val, exists := sel.Attr("datetime")
				if exists {
					item.Date = val
				}
			}
		}
		if req.ContentSelector != "" {
			sel := getSelection(req.ContentSelector)
			html, err := sel.Html()
			if err != nil {
				logrus.Infof("Warning: Failed to extract content using selector '%s' for item %d in feed %s: %v",
					req.ContentSelector, i, req.URL, err)
				item.Content = ""
			} else {
				item.Content = html
			}
		}

		items = append(items, item)
	})

	c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{
		StatusCode: 0,
		Data:       items,
	})
}
