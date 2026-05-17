package adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSimpleLLMCallRetriesDifferentModelsBeforeRepeating(t *testing.T) {
	llmClients = sync.Map{}
	t.Cleanup(func() {
		llmClients = sync.Map{}
	})

	var mu sync.Mutex
	requestedModels := make([]string, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
			http.Error(w, "unexpected path", http.StatusNotFound)
			return
		}

		var payload struct {
			Model string `json:"model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		mu.Lock()
		requestedModels = append(requestedModels, payload.Model)
		mu.Unlock()

		http.Error(w, "temporary upstream failure", http.StatusBadGateway)
	}))
	t.Cleanup(server.Close)

	t.Setenv("FC_LLM_API_TYPE", "openai")
	t.Setenv("FC_LLM_API_BASE", server.URL)
	t.Setenv("FC_LLM_API_KEY", "test-key")
	t.Setenv("FC_LLM_API_MODEL", "")

	_, err := simpleLLMCall("model-a, model-a, model-b", "hello", llmRetryConfig{
		attemptsPerModel: 2,
		delay:            time.Nanosecond,
		maxDelay:         time.Nanosecond,
	})
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.Len(t, requestedModels, 4)

	modelCounts := map[string]int{}
	for i, requestedModel := range requestedModels {
		modelCounts[requestedModel]++
		if i > 0 {
			require.NotEqual(t, requestedModels[i-1], requestedModel, "retry should try another configured model before repeating the same model")
		}
	}
	require.Equal(t, map[string]int{"model-a": 2, "model-b": 2}, modelCounts)
}

func TestSimpleLLMCallRetriesNextModelWithoutBackoff(t *testing.T) {
	llmClients = sync.Map{}
	t.Cleanup(func() {
		llmClients = sync.Map{}
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "temporary upstream failure", http.StatusBadGateway)
	}))
	t.Cleanup(server.Close)

	t.Setenv("FC_LLM_API_TYPE", "openai")
	t.Setenv("FC_LLM_API_BASE", server.URL)
	t.Setenv("FC_LLM_API_KEY", "test-key")
	t.Setenv("FC_LLM_API_MODEL", "")

	start := time.Now()
	_, err := simpleLLMCall("model-a,model-b", "hello", llmRetryConfig{
		attemptsPerModel: 1,
		delay:            100 * time.Millisecond,
		maxDelay:         200 * time.Millisecond,
	})
	require.Error(t, err)
	require.Less(t, time.Since(start), 150*time.Millisecond, "switching to another configured model should not wait for backoff")
}
