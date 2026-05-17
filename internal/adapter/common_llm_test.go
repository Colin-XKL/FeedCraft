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
	originalAttemptsPerModel := llmRetryAttemptsPerModel
	originalDelay := llmRetryDelay
	originalMaxDelay := llmRetryMaxDelay
	llmRetryAttemptsPerModel = 1
	llmRetryDelay = time.Nanosecond
	llmRetryMaxDelay = time.Nanosecond
	t.Cleanup(func() {
		llmRetryAttemptsPerModel = originalAttemptsPerModel
		llmRetryDelay = originalDelay
		llmRetryMaxDelay = originalMaxDelay
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

	_, err := SimpleLLMCall("model-a,model-b", "hello")
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.GreaterOrEqual(t, len(requestedModels), 2)
	require.NotEqual(t, requestedModels[0], requestedModels[1], "retry should try another configured model before repeating the same model")
}
