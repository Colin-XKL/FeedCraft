package config

import (
	"FeedCraft/internal/constant"
)

// --- Fetcher-specific Configurations ---

// HttpFetcherConfig holds the configuration for a simple HTTP GET fetcher.
type HttpFetcherConfig struct {
	URL            string            `json:"url"`
	Headers        map[string]string `json:"headers,omitempty"`
	UseBrowserless bool              `json:"use_browserless,omitempty"`
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
	ItemsIterator string `json:"items_iterator"`
	Title         string `json:"title"`
	Link          string `json:"link"`
	Date          string `json:"date,omitempty"`
	Description   string `json:"description,omitempty"`
	// ... other fields
}

// --- Feed-level Metadata Configuration ---

// FeedMetaConfig holds overrides for the final feed's metadata.
// These values will be used to replace any metadata parsed from the source.
type FeedMetaConfig struct {
	Title       string `json:"title,omitempty"`       // Force-override the feed title
	Link        string `json:"link,omitempty"`        // Force-override the feed's website link
	Description string `json:"description,omitempty"` // Force-override the feed description
	AuthorName  string `json:"author_name,omitempty"`
	AuthorEmail string `json:"author_email,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
}

// --- Top-level Source Configuration ---

// SourceConfig is the Go struct for the Recipe.SourceConfig field.
// It composes configurations for Fetcher, Parser, and Feed Metadata components.
type SourceConfig struct {
	Type constant.SourceType `json:"type"` // e.g., "rss", "html", "json"

	// Optional: Defines overrides for the final feed's metadata.
	FeedMeta *FeedMetaConfig `json:"feed_meta,omitempty"`

	// Fetcher configurations - only one should be non-nil for a given recipe.
	HttpFetcher *HttpFetcherConfig `json:"http_fetcher,omitempty"`
	// CurlFetcher *CurlFetcherConfig `json:"curl_fetcher,omitempty"` // Example for future use

	// Parser configurations - only one should be non-nil for a given recipe.
	// Note: RSS parsing doesn't require a specific config struct.
	HtmlParser *HtmlParserConfig `json:"html_parser,omitempty"`
	JsonParser *JsonParserConfig `json:"json_parser,omitempty"`
}
