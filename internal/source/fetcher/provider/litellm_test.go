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

func TestLiteLLMProvider_EndToEnd(t *testing.T) {
	mockResponse := `{
  "results": [
    {
      "url": "https://example.com/litellm1",
      "title": "LiteLLM Title 1",
      "snippet": "This is a snippet for the first result."
    },
    {
      "url": "https://example.com/litellm2",
      "title": "LiteLLM Title 2",
      "snippet": "Snippet for the second LiteLLM search result."
    }
  ]
}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, util.DefaultFeedUserAgent(), r.Header.Get("User-Agent"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	cfg := &config.SearchProviderConfig{
		APIUrl: ts.URL,
	}

	provider := NewLiteLLMProvider(cfg)
	data, err := provider.Fetch(context.Background(), "test litellm query")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	parserConfig := provider.GetDefaultParserConfig()
	jsonParser := &parser.JsonParser{Config: parserConfig}

	feed, err := jsonParser.Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	assert.Len(t, feed.Articles, 2)
	assert.Equal(t, "LiteLLM Title 1", feed.Articles[0].Title)
	assert.Equal(t, "https://example.com/litellm1", feed.Articles[0].Link)
	assert.Equal(t, "This is a snippet for the first result.", feed.Articles[0].Description)
	assert.Equal(t, "LiteLLM Title 2", feed.Articles[1].Title)
	assert.Equal(t, "https://example.com/litellm2", feed.Articles[1].Link)
	assert.Equal(t, "Snippet for the second LiteLLM search result.", feed.Articles[1].Description)
}
