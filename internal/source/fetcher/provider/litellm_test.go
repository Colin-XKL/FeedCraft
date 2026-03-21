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

func TestLiteLLMProvider_EndToEnd(t *testing.T) {
	// Standard API response from LiteLLM (often similar to OpenAI/Bing search results)
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
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	cfg := &config.SearchProviderConfig{
		APIUrl: ts.URL,
	}

	provider := NewLiteLLMProvider(cfg)

	// 1. Fetch the data
	data, err := provider.Fetch(context.Background(), "test litellm query")
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
	assert.Equal(t, "LiteLLM Title 1", feed.Items[0].Title)
	assert.Equal(t, "https://example.com/litellm1", feed.Items[0].Link)
	assert.Equal(t, "This is a snippet for the first result.", feed.Items[0].Description)

	// Second item verification
	assert.Equal(t, "LiteLLM Title 2", feed.Items[1].Title)
	assert.Equal(t, "https://example.com/litellm2", feed.Items[1].Link)
	assert.Equal(t, "Snippet for the second LiteLLM search result.", feed.Items[1].Description)
}
