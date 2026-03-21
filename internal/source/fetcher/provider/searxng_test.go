package provider

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/source/parser"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearXNGProvider_EndToEnd(t *testing.T) {
	// Standard API response from SearXNG
	mockResponse := `{
  "query": "test query",
  "number_of_results": 2,
  "results": [
    {
      "url": "https://example.com/1",
      "title": "Example Domain 1",
      "content": "This domain is for use in illustrative examples.",
      "publishedDate": "2023-10-01T10:00:00Z",
      "engine": "google"
    },
    {
      "url": "https://example.com/2",
      "title": "Example Domain 2",
      "content": "More illustrative examples here.",
      "engine": "bing"
    }
  ],
  "answers": [],
  "corrections": [],
  "infoboxes": [],
  "suggestions": []
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "test query", r.URL.Query().Get("q"))
		assert.Equal(t, "json", r.URL.Query().Get("format"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	cfg := &config.SearchProviderConfig{
		APIUrl: ts.URL,
		SearXNG: struct {
			Engines string `json:"engines"`
		}{},
	}

	provider := NewSearXNGProvider(cfg)

	// 1. Fetch the data
	data, err := provider.Fetch(context.Background(), "test query")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// 2. Parse the data using the default parser configuration
	parserConfig := provider.GetDefaultParserConfig()
	jsonParser := &parser.JsonParser{Config: parserConfig}

	feed, err := jsonParser.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, feed)

	// 3. Verify the parsed feed items
	assert.Len(t, feed.Items, 2)

	// First item verification
	assert.Equal(t, "Example Domain 1", feed.Items[0].Title)
	assert.Equal(t, "https://example.com/1", feed.Items[0].Link)
	assert.Equal(t, "This domain is for use in illustrative examples.", feed.Items[0].Description)
	assert.Equal(t, "2023-10-01T10:00:00Z", feed.Items[0].Published)

	// Second item verification (missing publishedDate)
	assert.Equal(t, "Example Domain 2", feed.Items[1].Title)
	assert.Equal(t, "https://example.com/2", feed.Items[1].Link)
	assert.Equal(t, "More illustrative examples here.", feed.Items[1].Description)
	assert.Empty(t, feed.Items[1].Published)
}
