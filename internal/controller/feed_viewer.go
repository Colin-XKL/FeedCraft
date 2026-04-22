package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/craft"
	"FeedCraft/internal/model"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
)

type FeedViewerPreviewReq struct {
	InputURL  string `json:"input_url" form:"input_url" binding:"required"`
	CraftName string `json:"craft_name" form:"craft_name"`
}

type FeedViewerPreview struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Link        string                  `json:"link"`
	FeedURL     string                  `json:"feedUrl"`
	Copyright   string                  `json:"copyright"`
	Image       *FeedViewerPreviewImage `json:"image,omitempty"`
	Items       []FeedViewerPreviewItem `json:"items"`
}

type FeedViewerPreviewImage struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type FeedViewerPreviewItem struct {
	GUID           string `json:"guid"`
	Title          string `json:"title"`
	Link           string `json:"link"`
	PubDate        string `json:"pubDate"`
	IsoDate        string `json:"isoDate"`
	Content        string `json:"content"`
	ContentSnippet string `json:"contentSnippet"`
}

func PreviewFeedViewer(c *gin.Context) {
	var req FeedViewerPreviewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: "Please enter a valid http(s) feed URL"})
		return
	}

	if err := validateFeedViewerURL(req.InputURL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	feed, err := loadFeedViewerPreview(c, req)
	if err != nil {
		status, msg := classifyFeedViewerError(err)
		c.JSON(status, util.APIResponse[any]{StatusCode: -1, Msg: msg})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[FeedViewerPreview]{
		StatusCode: 0,
		Data:       buildFeedViewerPreview(feed, req.InputURL),
	})
}

func loadFeedViewerPreview(c *gin.Context, req FeedViewerPreviewReq) (*model.CraftFeed, error) {
	cfg := &config.SourceConfig{
		Type: constant.SourceRSS,
		HttpFetcher: &config.HttpFetcherConfig{
			URL: req.InputURL,
		},
	}

	factory, err := source.Get(constant.SourceRSS)
	if err != nil {
		return nil, fmt.Errorf("factory not found: %w", err)
	}

	src, err := factory(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	feed, err := src.Fetch(c.Request.Context())
	if err != nil {
		return nil, err
	}

	if req.CraftName == "" || req.CraftName == "proxy" {
		return feed, nil
	}

	craftedFeed, err := buildCraftPreview(feed, req.InputURL, req.CraftName)
	if err != nil {
		return nil, err
	}

	return craftedFeed, nil
}

func buildCraftPreview(feed *model.CraftFeed, inputURL, craftName string) (*model.CraftFeed, error) {
	atomXML, err := feed.ToFeedsFeed().ToAtom()
	if err != nil {
		return nil, err
	}

	parsedFeed, err := gofeed.NewParser().ParseString(atomXML)
	if err != nil {
		return nil, err
	}

	craftedFeed, err := craft.ProcessFeed(parsedFeed, inputURL, craftName)
	if err != nil {
		return nil, err
	}

	return model.FromFeedsFeed(craftedFeed), nil
}

func buildFeedViewerPreview(feed *model.CraftFeed, inputURL string) FeedViewerPreview {
	preview := FeedViewerPreview{
		Title:       feed.Title,
		Description: feed.Description,
		Link:        feed.Link,
		FeedURL:     inputURL,
		Copyright:   feed.Copyright,
		Items:       make([]FeedViewerPreviewItem, 0, len(feed.Articles)),
	}

	if feed.ImageURL != "" || feed.ImageTitle != "" {
		preview.Image = &FeedViewerPreviewImage{
			URL:   feed.ImageURL,
			Title: feed.ImageTitle,
		}
	}

	for _, article := range feed.Articles {
		if article == nil {
			continue
		}

		contentSnippet := article.Description
		if contentSnippet == "" {
			contentSnippet = article.Content
		}

		preview.Items = append(preview.Items, FeedViewerPreviewItem{
			GUID:           article.Id,
			Title:          article.Title,
			Link:           article.Link,
			PubDate:        formatFeedViewerTime(article.Created, article.Updated),
			IsoDate:        formatFeedViewerISOTime(article.Created, article.Updated),
			Content:        article.Content,
			ContentSnippet: contentSnippet,
		})
	}

	return preview
}

func formatFeedViewerTime(primary, fallback time.Time) string {
	if !primary.IsZero() {
		return primary.Format(time.RFC1123Z)
	}
	if !fallback.IsZero() {
		return fallback.Format(time.RFC1123Z)
	}
	return ""
}

func formatFeedViewerISOTime(primary, fallback time.Time) string {
	if !primary.IsZero() {
		return primary.UTC().Format(time.RFC3339)
	}
	if !fallback.IsZero() {
		return fallback.UTC().Format(time.RFC3339)
	}
	return ""
}

func validateFeedViewerURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL == nil {
		return errors.New("Please enter a valid http(s) feed URL")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("Please enter a valid http(s) feed URL")
	}
	if parsedURL.Hostname() == "" {
		return errors.New("Please enter a valid http(s) feed URL")
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return fmt.Errorf("Unable to resolve this URL: %w", err)
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return fmt.Errorf("Access to private IP %s is forbidden", ip.String())
		}
	}

	return nil
}

func classifyFeedViewerError(err error) (int, string) {
	msg := err.Error()
	if strings.HasPrefix(msg, "all items failed to process. last error: ") {
		msg = strings.TrimPrefix(msg, "all items failed to process. last error: ")
	}

	switch {
	case strings.Contains(msg, "browserless service returned status"):
		return http.StatusOK, humanizeBrowserlessStatus(msg)
	case strings.Contains(msg, "http status not ok:"):
		return http.StatusOK, humanizeFeedViewerHTTPStatus(msg)
	case strings.Contains(msg, "http get failed:"), strings.Contains(msg, "browserless fetch failed:"), strings.Contains(msg, "failed to read response body:"), strings.Contains(msg, "Unable to resolve this URL"):
		return http.StatusOK, "Unable to fetch this URL. Please check the address and try again."
	case strings.Contains(msg, "parse failed:"), strings.Contains(msg, "invalid XML"):
		return http.StatusOK, "The URL is reachable, but it does not appear to be a valid RSS or Atom feed."
	case strings.Contains(msg, "not a valid craft name"):
		return http.StatusBadRequest, "Please select a valid craft before comparing feeds."
	default:
		return http.StatusInternalServerError, "Failed to preview this feed due to an internal error."
	}
}

func humanizeBrowserlessStatus(msg string) string {
	status := strings.TrimSpace(strings.TrimPrefix(msg, "browserless service returned status"))
	if status == "" {
		return "Browserless service failed to render the URL."
	}
	return fmt.Sprintf("Browserless service failed to render the URL (returned status %s). Please check the address or the browserless service.", status)
}

func humanizeFeedViewerHTTPStatus(msg string) string {
	status := strings.TrimSpace(strings.TrimPrefix(msg, "fetch failed: http status not ok:"))
	if status == "" {
		status = strings.TrimSpace(strings.TrimPrefix(msg, "http status not ok:"))
	}
	if status == "" {
		return "Unable to fetch this URL. Please check the address and try again."
	}
	return fmt.Sprintf("The source returned %s, so the feed could not be loaded.", status)
}
