package provider

import (
	"FeedCraft/internal/config"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteLLMProvider_Fetch(t *testing.T) {
	mockResponse := `{"results": [{"title": "Test Title", "url": "https://example.com", "snippet": "Test Content"}]}`
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
	data, err := provider.Fetch(context.Background(), "test")

	assert.NoError(t, err)
	assert.Contains(t, string(data), "Test Title")
}

func TestLiteLLMProvider_GetDefaultParserConfig(t *testing.T) {
	provider := NewLiteLLMProvider(&config.SearchProviderConfig{})
	cfg := provider.GetDefaultParserConfig()

	assert.Equal(t, ".results", cfg.ItemsIterator)
	assert.Equal(t, ".title", cfg.Title)
	assert.Equal(t, ".url", cfg.Link)
	assert.Equal(t, ".snippet", cfg.Description)
}
