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
)

type FetchReq struct {
	URL string `json:"url" binding:"required"`
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

func FetchURL(c *gin.Context) {
	var req FetchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	if err := validateURL(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	// Try standard HTTP request first (simulating a browser user agent)
	client := resty.New()
	client.SetTimeout(craft.DefaultExtractFulltextTimeout)
	resp, err := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
		Get(req.URL)

	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Fetch failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[string]{
		StatusCode: 0,
		Data:       resp.String(),
	})
}

func ParseRSS(c *gin.Context) {
	var req ParseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	var htmlContent string
	if req.HTML != "" {
		htmlContent = req.HTML
	} else if req.URL != "" {
		if err := validateURL(req.URL); err != nil {
			c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
			return
		}
		// Fetch if HTML not provided
		client := resty.New()
		client.SetTimeout(craft.DefaultExtractFulltextTimeout)
		resp, err := client.R().
			SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36").
			Get(req.URL)
		if err != nil {
			c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Fetch failed: " + err.Error()})
			return
		}
		htmlContent = resp.String()
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
		if req.TitleSelector != "" {
			item.Title = strings.TrimSpace(s.Find(req.TitleSelector).Text())
		}
		if req.LinkSelector != "" {
			// Try to get href
			href, exists := s.Find(req.LinkSelector).Attr("href")
			if exists {
				item.Link = href
			} else {
				item.Link = strings.TrimSpace(s.Find(req.LinkSelector).Text())
			}
		}
		if req.DateSelector != "" {
			item.Date = strings.TrimSpace(s.Find(req.DateSelector).Text())
			// Check datetime attr if text is empty
			if item.Date == "" {
				val, exists := s.Find(req.DateSelector).Attr("datetime")
				if exists {
					item.Date = val
				}
			}
		}
		if req.ContentSelector != "" {
			html, _ := s.Find(req.ContentSelector).Html()
			item.Content = html // Return HTML for content
		}

		items = append(items, item)
	})

	c.JSON(http.StatusOK, util.APIResponse[[]ParsedItem]{
		StatusCode: 0,
		Data:       items,
	})
}
