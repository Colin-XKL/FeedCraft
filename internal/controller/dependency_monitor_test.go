package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckEmbeddingFallsBackToLLMEndpointForConfiguredStatus(t *testing.T) {
	t.Setenv("FC_LLM_API_TYPE", "openai")
	t.Setenv("FC_LLM_API_BASE", "https://llm.example/v1")
	t.Setenv("FC_LLM_API_KEY", "llm-secret")
	t.Setenv("FC_LLM_API_MODEL", "gpt-4o")
	t.Setenv("FC_EMBEDDING_API_TYPE", "")
	t.Setenv("FC_EMBEDDING_API_BASE", "")
	t.Setenv("FC_EMBEDDING_API_KEY", "")
	t.Setenv("FC_EMBEDDING_API_MODEL", "text-embedding-3-small")

	status := checkEmbedding(false)

	if status.Name != "Embedding Service" {
		t.Fatalf("name = %q", status.Name)
	}
	if status.Status != "Configured" {
		t.Fatalf("status = %q", status.Status)
	}
	if status.Error != "" {
		t.Fatalf("unexpected error = %q", status.Error)
	}
	for _, want := range []string{"Type: openai", "Model: text-embedding-3-small", "Base: https://llm.example/v1", "Key: llm-******cret"} {
		if !strings.Contains(status.Details, want) {
			t.Fatalf("details %q does not contain %q", status.Details, want)
		}
	}
}

func TestCheckEmbeddingActiveCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embeddings" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("unexpected Authorization header %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"data": []map[string]any{
				{
					"object":    "embedding",
					"index":     0,
					"embedding": []float64{1, 0, 0},
				},
			},
			"model": "test-embedding",
		})
	}))
	defer server.Close()

	t.Setenv("FC_EMBEDDING_API_TYPE", "openai")
	t.Setenv("FC_EMBEDDING_API_BASE", server.URL)
	t.Setenv("FC_EMBEDDING_API_KEY", "test-key")
	t.Setenv("FC_EMBEDDING_API_MODEL", "test-embedding")

	status := checkEmbedding(true)

	if status.Status != "Healthy" {
		t.Fatalf("status = %q, error = %q", status.Status, status.Error)
	}
	if status.Latency == "" {
		t.Fatal("expected latency to be populated")
	}
}

func TestCheckBrowserProviderUsesCloakBrowserVersionEndpoint(t *testing.T) {
	var requestedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		assert.Equal(t, "feedcraft", r.URL.Query().Get("fingerprint"))
		_, _ = w.Write([]byte(`{"webSocketDebuggerUrl":"ws://example/devtools/browser/1"}`))
	}))
	defer server.Close()

	env := viper.New()
	env.Set("BROWSER_PROVIDER", "cdp")
	env.Set("BROWSER_ENDPOINT", server.URL+"?fingerprint=feedcraft")

	status := checkBrowserless(env, true)

	assert.Equal(t, "Browser Provider", status.Name)
	assert.Equal(t, "Healthy", status.Status)
	assert.Equal(t, "/json/version", requestedPath)
	assert.Contains(t, status.Details, "Provider: cdp")
	assert.Contains(t, status.Details, server.URL)
}

func TestCheckBrowserProviderRejectsUnsupportedProvider(t *testing.T) {
	env := viper.New()
	env.Set("BROWSER_PROVIDER", "cdp-typo")
	env.Set("BROWSER_ENDPOINT", "http://example.com")

	status := checkBrowserless(env, true)

	assert.Equal(t, "Browser Provider", status.Name)
	assert.Equal(t, "Unhealthy", status.Status)
	assert.Contains(t, status.Error, "unsupported browser provider")
}
