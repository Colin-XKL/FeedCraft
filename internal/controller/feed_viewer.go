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

type EmbeddingFilterPreviewReq struct {
	InputURL         string   `json:"input_url" binding:"required"`
	Anchors          string   `json:"anchors" binding:"required"`
	Threshold        *float64 `json:"threshold"`
	Mode             string   `json:"mode"`
	MaxContentLength *int     `json:"max_content_length"`
	Instruction      string   `json:"instruction"`
}

type embeddingFilterPreviewConfig struct {
	inputURL         string
	anchors          []string
	threshold        float64
	mode             craft.EmbeddingFilterMode
	maxContentLength int
	instruction      string
}

const (
	maxEmbeddingFilterPreviewAnchors           = 20
	maxEmbeddingFilterPreviewAnchorLength      = 500
	maxEmbeddingFilterPreviewInstructionLength = 1000
	maxEmbeddingFilterPreviewContentLength     = 8000
	maxEmbeddingFilterPreviewItems             = 80
)

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

func PreviewEmbeddingFilter(c *gin.Context) {
	var req EmbeddingFilterPreviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	cfg, err := normalizeEmbeddingFilterPreviewRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}
	if err := validateFeedViewerURL(cfg.inputURL); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: formatFeedViewerValidationError(err)})
		return
	}

	feed, err := loadFeedViewerPreview(c, FeedViewerPreviewReq{InputURL: cfg.inputURL})
	if err != nil {
		status, msg := classifyFeedViewerError(err)
		c.JSON(status, util.APIResponse[any]{StatusCode: -1, Msg: msg})
		return
	}
	if len(feed.Articles) > maxEmbeddingFilterPreviewItems {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{
			StatusCode: -1,
			Msg:        fmt.Sprintf("embedding filter preview supports at most %d feed items", maxEmbeddingFilterPreviewItems),
		})
		return
	}

	craftedFeed, err := buildCraftPreviewWithOptions(feed, cfg.inputURL, []craft.LegacyCraftOption{craft.OptionEmbeddingFilterWithContext(
		c.Request.Context(),
		cfg.anchors,
		cfg.threshold,
		cfg.maxContentLength,
		cfg.instruction,
		cfg.mode,
	)})
	if err != nil {
		status, msg := classifyFeedViewerError(err)
		c.JSON(status, util.APIResponse[any]{StatusCode: -1, Msg: msg})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[FeedViewerPreview]{
		StatusCode: 0,
		Data:       buildFeedViewerPreview(craftedFeed, cfg.inputURL),
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

func normalizeEmbeddingFilterPreviewRequest(req EmbeddingFilterPreviewReq) (embeddingFilterPreviewConfig, error) {
	cfg := embeddingFilterPreviewConfig{
		inputURL:         strings.TrimSpace(req.InputURL),
		threshold:        0.6,
		mode:             craft.EmbeddingIncludeMode,
		maxContentLength: 2000,
		instruction:      strings.TrimSpace(req.Instruction),
	}

	for _, rawAnchor := range strings.Split(req.Anchors, "\n") {
		anchor := strings.TrimSpace(rawAnchor)
		if anchor != "" {
			if len([]rune(anchor)) > maxEmbeddingFilterPreviewAnchorLength {
				return cfg, fmt.Errorf("anchors must be at most %d characters each", maxEmbeddingFilterPreviewAnchorLength)
			}
			cfg.anchors = append(cfg.anchors, anchor)
		}
	}
	if len(cfg.anchors) == 0 {
		return cfg, errors.New("anchors parameter is required")
	}
	if len(cfg.anchors) > maxEmbeddingFilterPreviewAnchors {
		return cfg, fmt.Errorf("anchors supports at most %d entries", maxEmbeddingFilterPreviewAnchors)
	}
	if len([]rune(cfg.instruction)) > maxEmbeddingFilterPreviewInstructionLength {
		return cfg, fmt.Errorf("instruction must be at most %d characters", maxEmbeddingFilterPreviewInstructionLength)
	}

	if req.Threshold != nil {
		if *req.Threshold < 0 || *req.Threshold > 1 {
			return cfg, errors.New("threshold must be between 0 and 1")
		}
		cfg.threshold = *req.Threshold
	}

	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	switch mode {
	case "", "include":
		cfg.mode = craft.EmbeddingIncludeMode
	case "exclude":
		cfg.mode = craft.EmbeddingExcludeMode
	default:
		return cfg, errors.New("mode must be include or exclude")
	}

	if req.MaxContentLength != nil {
		if *req.MaxContentLength <= 0 {
			return cfg, errors.New("max_content_length must be greater than 0")
		}
		if *req.MaxContentLength > maxEmbeddingFilterPreviewContentLength {
			return cfg, fmt.Errorf("max_content_length must be at most %d", maxEmbeddingFilterPreviewContentLength)
		}
		cfg.maxContentLength = *req.MaxContentLength
	}

	return cfg, nil
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

func buildCraftPreviewWithOptions(feed *model.CraftFeed, inputURL string, options []craft.LegacyCraftOption) (*model.CraftFeed, error) {
	atomXML, err := feed.ToFeedsFeed().ToAtom()
	if err != nil {
		return nil, err
	}

	parsedFeed, err := gofeed.NewParser().ParseString(atomXML)
	if err != nil {
		return nil, err
	}

	craftedFeed, err := craft.NewCraftedFeedFromGofeed(parsedFeed, inputURL, options...)
	if err != nil {
		return nil, err
	}

	return model.FromFeedsFeed(craftedFeed.OutputFeed), nil
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
	case isEmbeddingFilterConfigurationError(msg):
		return http.StatusBadRequest, humanizeEmbeddingFilterError(msg)
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

func isEmbeddingFilterConfigurationError(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	if !strings.Contains(lowerMsg, "[embedding-filter]") {
		return false
	}
	return strings.Contains(lowerMsg, "failed to load embedding config") ||
		strings.Contains(lowerMsg, "fc_embedding_") ||
		strings.Contains(lowerMsg, "anchors parameter is required")
}

func humanizeEmbeddingFilterError(msg string) string {
	detail := strings.TrimSpace(msg)
	prefixes := []string{
		"[embedding-filter] failed to compute anchor vectors:",
		"failed to compute anchor vectors:",
		"failed to load embedding config:",
		"[embedding-filter]",
	}
	for {
		trimmed := detail
		for _, prefix := range prefixes {
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		}
		if trimmed == detail {
			break
		}
		detail = trimmed
	}
	if detail == "" {
		return "Embedding filter is not configured correctly."
	}
	return "Embedding filter is not configured correctly: " + detail
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
