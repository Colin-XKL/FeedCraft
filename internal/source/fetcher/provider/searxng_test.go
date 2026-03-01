package provider

import (
	"FeedCraft/internal/config"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearXNGProvider_Fetch(t *testing.T) {
	// Create a mock SearXNG server
	mockResponse := `{"query": "test", "results": [{"title": "Test Title", "url": "https://example.com", "content": "Test Content", "publishedDate": "2023-10-01"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "test", r.URL.Query().Get("q"))
		assert.Equal(t, "json", r.URL.Query().Get("format"))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	cfg := &config.SearchProviderConfig{
		APIUrl: ts.URL,
		SearXNG: struct {
			Engines string `json:"engines"`
		}{},
	}

	provider := NewSearXNGProvider(cfg)
	data, err := provider.Fetch(context.Background(), "test")

	assert.NoError(t, err)
	assert.Contains(t, string(data), "Test Title")
}

func TestSearXNGProvider_GetDefaultParserConfig(t *testing.T) {
	provider := NewSearXNGProvider(&config.SearchProviderConfig{})
	cfg := provider.GetDefaultParserConfig()

	assert.Equal(t, ".results", cfg.ItemsIterator)
	assert.Equal(t, ".title", cfg.Title)
	assert.Equal(t, ".url", cfg.Link)
	assert.Equal(t, ".content", cfg.Description)
	assert.Equal(t, ".publishedDate", cfg.Date)
}
