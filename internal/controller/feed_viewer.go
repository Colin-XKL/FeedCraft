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

const feedViewerInvalidURLMessage = "Please enter a valid http(s) feed URL"
const feedViewerInvalidURLError = "please enter a valid http(s) feed URL"

func PreviewFeedViewer(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, util.APIResponse[any]{StatusCode: -1, Msg: "Only GET requests are allowed for feed preview"})
		return
	}

	var req FeedViewerPreviewReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: feedViewerInvalidURLMessage})
		return
	}

	if err := validateFeedViewerURL(req.InputURL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: formatFeedViewerValidationError(err)})
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

func formatFeedViewerValidationError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if msg == "" {
		return ""
	}
	return strings.ToUpper(msg[:1]) + msg[1:]
}

func validateFeedViewerURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL == nil {
		return errors.New(feedViewerInvalidURLError)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New(feedViewerInvalidURLError)
	}
	if parsedURL.Hostname() == "" {
		return errors.New(feedViewerInvalidURLError)
	}

	return nil
}

func classifyFeedViewerError(err error) (int, string) {
	msg := err.Error()
	msg = strings.TrimPrefix(msg, "all items failed to process. last error: ")
	lowerMsg := strings.ToLower(msg)

	switch {
	case strings.Contains(lowerMsg, "browserless service returned status"):
		return http.StatusOK, humanizeBrowserlessStatus(msg)
	case strings.Contains(lowerMsg, "browser cdp render failed"),
		strings.Contains(lowerMsg, "browser provider not configured"),
		strings.Contains(lowerMsg, "browser cdp service returned status"),
		strings.Contains(lowerMsg, "browser cdp version request failed"),
		strings.Contains(lowerMsg, "failed to decode browser cdp version response"),
		strings.Contains(lowerMsg, "browser cdp version response missing websocketdebuggerurl"),
		strings.Contains(lowerMsg, "unsupported browser provider"):
		return http.StatusOK, "Browser provider failed to render the URL. Please check the address or the browser provider service."
	case strings.Contains(lowerMsg, "http status not ok:"):
		return http.StatusOK, humanizeFeedViewerHTTPStatus(msg)
	case strings.Contains(lowerMsg, "http get failed:"), strings.Contains(lowerMsg, "browserless fetch failed:"), strings.Contains(lowerMsg, "failed to read response body:"), strings.Contains(lowerMsg, "unable to resolve this url"):
		return http.StatusOK, "Unable to fetch this URL. Please check the address and try again."
	case strings.Contains(lowerMsg, "parse failed:"), strings.Contains(lowerMsg, "invalid xml"):
		return http.StatusOK, "The URL is reachable, but it does not appear to be a valid RSS, Atom, or JSON feed."
	case strings.Contains(lowerMsg, "not a valid craft name"):
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
