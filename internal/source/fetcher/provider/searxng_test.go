package provider

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/source/parser"
	"FeedCraft/internal/util"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearXNGProvider_EndToEnd(t *testing.T) {
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
		assert.Equal(t, util.DefaultFeedUserAgent(), r.Header.Get("User-Agent"))
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
	data, err := provider.Fetch(context.Background(), "test query")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	parserConfig := provider.GetDefaultParserConfig()
	jsonParser := &parser.JsonParser{Config: parserConfig}

	feed, err := jsonParser.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	assert.Len(t, feed.Articles, 2)
	assert.Equal(t, "Example Domain 1", feed.Articles[0].Title)
	assert.Equal(t, "https://example.com/1", feed.Articles[0].Link)
	assert.Equal(t, "This domain is for use in illustrative examples.", feed.Articles[0].Description)
	assert.False(t, feed.Articles[0].Created.IsZero())
	assert.Equal(t, "Example Domain 2", feed.Articles[1].Title)
	assert.Equal(t, "https://example.com/2", feed.Articles[1].Link)
	assert.Equal(t, "More illustrative examples here.", feed.Articles[1].Description)
	assert.True(t, feed.Articles[1].Created.IsZero())
}
