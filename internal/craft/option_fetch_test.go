package craft

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCraftedFeedFromUrlUsesFeedFetcherUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Feed</title>
    <link>https://example.com/</link>
    <description>Test feed</description>
    <item>
      <title>Item 1</title>
      <link>https://example.com/item-1</link>
      <description>Hello</description>
    </item>
  </channel>
</rss>`))
	}))
	defer server.Close()

	crafted, err := NewCraftedFeedFromUrl(server.URL)
	if err != nil {
		t.Fatalf("NewCraftedFeedFromUrl returned error: %v", err)
	}
	if crafted.OutputFeed == nil {
		t.Fatal("expected output feed")
	}
	if gotUA != "FeedCraft/2.0" {
		t.Fatalf("expected feed fetcher user agent, got %q", gotUA)
	}
}
