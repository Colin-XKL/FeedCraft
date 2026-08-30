package config

import (
	"FeedCraft/internal/constant"
	"time"
)

const (
	HttpFetcherPurposeFeed = "feed"
	HttpFetcherPurposeHTML = "html"
)

const (
	BrowserNavigationActionClick           = "click"
	BrowserNavigationActionWaitForSelector = "wait_for_selector"
	BrowserNavigationActionWait            = "wait"

	FeedIconSourceAuto           = "auto"
	FeedIconSourceFaviconService = "favicon_service"
)

// BrowserNavigationAction describes a browser interaction to run after the page
// loads and before rendered HTML is captured.
type BrowserNavigationAction struct {
	Type       string `json:"type"`
	Selector   string `json:"selector,omitempty"`
	TimeoutMs  int64  `json:"timeout_ms,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
}

// --- Fetcher-specific Configurations ---

// HttpFetcherConfig holds the configuration for an HTTP fetcher.
// It supports some HTTP methods, headers, and request bodies.
type HttpFetcherConfig struct {
	URL               string                    `json:"url"`
	Method            string                    `json:"method,omitempty"`
	Headers           map[string]string         `json:"headers,omitempty"`
	Body              string                    `json:"body,omitempty"`
	UseBrowserless    bool                      `json:"use_browserless,omitempty"`
	NavigationActions []BrowserNavigationAction `json:"navigation_actions,omitempty"`
	Purpose           string                    `json:"purpose,omitempty"`
	HttpClientTimeout time.Duration             `json:"http_client_timeout,omitempty"`
}

// SearchFetcherConfig holds the configuration for search-based fetching.
type SearchFetcherConfig struct {
	Query        string `json:"query"`
	EnhancedMode bool   `json:"enhanced_mode,omitempty"`
}

// --- Parser-specific Configurations ---

// HtmlParserConfig holds the configuration for parsing HTML content.
type HtmlParserConfig struct {
	ItemSelector string `json:"item_selector"`
	Title        string `json:"title"`
	Link         string `json:"link"`
	Date         string `json:"date,omitempty"`
	Description  string `json:"description,omitempty"`
	Content      string `json:"content,omitempty"`
	// ... other field selectors
}

// JsonParserConfig holds the configuration for parsing JSON content.
type JsonParserConfig struct {
	ItemsIterator       string `json:"items_iterator"`
	Title               string `json:"title"`
	TitleTemplate       string `json:"title_template,omitempty"`
	Link                string `json:"link"`
	LinkTemplate        string `json:"link_template,omitempty"`
	Date                string `json:"date,omitempty"`
	DateTemplate        string `json:"date_template,omitempty"`
	Description         string `json:"description,omitempty"`
	DescriptionTemplate string `json:"description_template,omitempty"`
	// ... other fields
}

// WebMonitorParserConfig holds the configuration for monitoring changes in HTML content.
type WebMonitorParserConfig struct {
	Extractors          map[string]string `json:"extractors"`
	KeyFields           []string          `json:"key_fields"`
	TitleTemplate       string            `json:"title_template,omitempty"`
	DescriptionTemplate string            `json:"description_template,omitempty"`
	ContentTemplate     string            `json:"content_template,omitempty"`
}

// InboxSourceConfig holds the configuration for an inbox source.
type InboxSourceConfig struct {
	InboxID string `json:"inbox_id"`
}

// --- Feed-level Metadata Configuration ---

// FeedMetaConfig holds overrides for the final feed's metadata.
// These values will be used to replace any metadata parsed from the source.
type FeedMetaConfig struct {
	Title           string `json:"title,omitempty"`       // Force-override the feed title
	Link            string `json:"link,omitempty"`        // Force-override the feed's website link
	Description     string `json:"description,omitempty"` // Force-override the feed description
	AuthorName      string `json:"author_name,omitempty"`
	AuthorEmail     string `json:"author_email,omitempty"`
	Copyright       string `json:"copyright,omitempty"`
	IconSource      string `json:"icon_source,omitempty"`
	FaviconProvider string `json:"favicon_provider,omitempty"`
}

// --- Top-level Source Configuration ---

// SourceConfig is the Go struct for the Recipe.SourceConfig field.
// It composes configurations for Fetcher, Parser, and Feed Metadata components.
type SourceConfig struct {
	Type constant.SourceType `json:"type"` // e.g., "rss", "html", "json"

	// Optional: Defines overrides for the final feed's metadata.
	FeedMeta *FeedMetaConfig `json:"feed_meta,omitempty"`

	// Fetcher configurations - only one should be non-nil for a given recipe.
	HttpFetcher   *HttpFetcherConfig   `json:"http_fetcher,omitempty"`
	SearchFetcher *SearchFetcherConfig `json:"search_fetcher,omitempty"`
	// CurlFetcher *CurlFetcherConfig `json:"curl_fetcher,omitempty"` // Example for future use

	// Parser configurations - only one should be non-nil for a given recipe.
	// Note: RSS parsing doesn't require a specific config struct.
	HtmlParser       *HtmlParserConfig       `json:"html_parser,omitempty"`
	JsonParser       *JsonParserConfig       `json:"json_parser,omitempty"`
	WebMonitorParser *WebMonitorParserConfig `json:"web_monitor_parser,omitempty"`

	InboxSource *InboxSourceConfig `json:"inbox_source,omitempty"`
}
