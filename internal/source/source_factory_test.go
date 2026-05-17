package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"testing"
)

func TestRSSSourceFactorySetsFeedPurposeByDefault(t *testing.T) {
	cfg := &config.SourceConfig{
		Type: constant.SourceRSS,
		HttpFetcher: &config.HttpFetcherConfig{
			URL: "https://example.com/feed.xml",
		},
	}

	_, err := rssSourceFactory(cfg)
	if err != nil {
		t.Fatalf("rssSourceFactory returned error: %v", err)
	}
	if cfg.HttpFetcher.Purpose != config.HttpFetcherPurposeFeed {
		t.Fatalf("expected purpose %q, got %q", config.HttpFetcherPurposeFeed, cfg.HttpFetcher.Purpose)
	}
}

func TestHTMLSourceFactorySetsHTMLPurposeByDefault(t *testing.T) {
	cfg := &config.SourceConfig{
		Type: constant.SourceHTML,
		HttpFetcher: &config.HttpFetcherConfig{
			URL: "https://example.com/page",
		},
		HtmlParser: &config.HtmlParserConfig{
			ItemSelector: ".item",
			Title:        ".title",
			Link:         ".link",
		},
	}

	_, err := htmlSourceFactory(cfg)
	if err != nil {
		t.Fatalf("htmlSourceFactory returned error: %v", err)
	}
	if cfg.HttpFetcher.Purpose != config.HttpFetcherPurposeHTML {
		t.Fatalf("expected purpose %q, got %q", config.HttpFetcherPurposeHTML, cfg.HttpFetcher.Purpose)
	}
}

func TestJSONSourceFactorySetsHTMLPurposeByDefault(t *testing.T) {
	cfg := &config.SourceConfig{
		Type: constant.SourceJSON,
		HttpFetcher: &config.HttpFetcherConfig{
			URL: "https://example.com/api/items",
		},
		JsonParser: &config.JsonParserConfig{
			ItemsIterator: ".items",
			Title:         ".title",
			Link:          ".link",
		},
	}

	_, err := jsonSourceFactory(cfg)
	if err != nil {
		t.Fatalf("jsonSourceFactory returned error: %v", err)
	}
	if cfg.HttpFetcher.Purpose != config.HttpFetcherPurposeHTML {
		t.Fatalf("expected purpose %q, got %q", config.HttpFetcherPurposeHTML, cfg.HttpFetcher.Purpose)
	}
}
