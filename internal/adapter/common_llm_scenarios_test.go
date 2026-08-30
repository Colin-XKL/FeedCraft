package adapter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// chatCompletionRequest 描述 OpenAI 兼容 /chat/completions 请求体中我们关心的字段。
type chatCompletionRequest struct {
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	Messages    []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

// fastRetry 是测试用的极速重试配置，避免真实退避拖慢测试。
var fastRetry = llmRetryConfig{
	attemptsPerModel: 3,
	delay:            time.Nanosecond,
	maxDelay:         time.Nanosecond,
}

// resetLLMClients 清空全局客户端缓存，避免测试间互相污染。
func resetLLMClients(t *testing.T) {
	t.Helper()
	llmClients = sync.Map{}
	t.Cleanup(func() {
		llmClients = sync.Map{}
	})
}

// setOpenAILLMEnv 设置标准的 OpenAI 兼容环境变量，并清空 legacy 变量。
func setOpenAILLMEnv(t *testing.T, baseURL string) {
	t.Helper()
	t.Setenv("FC_LLM_API_TYPE", "openai")
	t.Setenv("FC_LLM_API_BASE", baseURL)
	t.Setenv("FC_LLM_API_KEY", "test-key")
	t.Setenv("FC_LLM_API_MODEL", "")
	t.Setenv("FC_OPENAI_ENDPOINT", "")
	t.Setenv("FC_OPENAI_AUTH_KEY", "")
	t.Setenv("FC_OPENAI_DEFAULT_MODEL", "")
}

// writeChatCompletion 写出一个最简单的成功 chat/completions 响应。
func writeChatCompletion(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		"choices": []map[string]any{
			{"message": map[string]any{"role": "assistant", "content": content}},
		},
		"usage": map[string]any{"total_tokens": 1},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// --- OpenAI 兼容 API 基础功能 ---

func TestSimpleLLMCallSuccessSendsExpectedRequest(t *testing.T) {
	resetLLMClients(t)

	var captured chatCompletionRequest
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/chat/completions", r.URL.Path)
		authHeader = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, "hello world")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	result, err := SimpleLLMCall("gpt-test", "please answer")

	require.NoError(t, err)
	require.Equal(t, "hello world", result)
	require.Equal(t, "Bearer test-key", authHeader, "OpenAI 兼容 API 必须发送 Bearer 鉴权头")
	require.Equal(t, "gpt-test", captured.Model)
	require.Len(t, captured.Messages, 1)
	require.Equal(t, "user", captured.Messages[0].Role)
	require.Equal(t, "please answer", captured.Messages[0].Content)
}

func TestSimpleLLMCallUsesDefaultModelFromEnv(t *testing.T) {
	resetLLMClients(t)

	var captured chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, "ok")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "env-default-model")

	result, err := SimpleLLMCall(UseDefaultModel, "hi")

	require.NoError(t, err)
	require.Equal(t, "ok", result)
	require.Equal(t, "env-default-model", captured.Model, "model 入参为空时应使用 FC_LLM_API_MODEL")
}

func TestSimpleLLMCallFallsBackToLegacyEnv(t *testing.T) {
	resetLLMClients(t)

	var captured chatCompletionRequest
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, "legacy ok")
	}))
	t.Cleanup(server.Close)

	// 只设置 legacy 变量，新变量留空，验证向后兼容。
	t.Setenv("FC_LLM_API_TYPE", "")
	t.Setenv("FC_LLM_API_BASE", "")
	t.Setenv("FC_LLM_API_KEY", "")
	t.Setenv("FC_LLM_API_MODEL", "")
	t.Setenv("FC_OPENAI_ENDPOINT", server.URL)
	t.Setenv("FC_OPENAI_AUTH_KEY", "legacy-key")
	t.Setenv("FC_OPENAI_DEFAULT_MODEL", "legacy-model")

	result, err := SimpleLLMCall(UseDefaultModel, "hi")

	require.NoError(t, err)
	require.Equal(t, "legacy ok", result)
	require.Equal(t, "Bearer legacy-key", authHeader)
	require.Equal(t, "legacy-model", captured.Model)
}

func TestSimpleLLMCallNoModelConfigured(t *testing.T) {
	resetLLMClients(t)
	setOpenAILLMEnv(t, "http://127.0.0.1:0")

	_, err := simpleLLMCall(UseDefaultModel, "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
}

func TestSimpleLLMCallOllamaRequiresBase(t *testing.T) {
	resetLLMClients(t)
	t.Setenv("FC_LLM_API_TYPE", "ollama")
	t.Setenv("FC_LLM_API_BASE", "")
	t.Setenv("FC_LLM_API_KEY", "")
	t.Setenv("FC_LLM_API_MODEL", "")
	t.Setenv("FC_OPENAI_ENDPOINT", "")

	_, err := simpleLLMCall("llama3", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "FC_LLM_API_BASE")
}

// --- 异常 / 边界场景 ---

func TestSimpleLLMCallEmptyChoicesReturnsError(t *testing.T) {
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[],"usage":{"total_tokens":0}}`))
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
	// 空 choices 被视为可重试错误，应重试 attemptsPerModel 次。
	require.Equal(t, int32(fastRetry.attemptsPerModel), atomic.LoadInt32(&calls))
}

func TestSimpleLLMCallNetworkInterruption(t *testing.T) {
	resetLLMClients(t)

	// 启动后立即关闭服务器，模拟网络中断 / 连接被拒。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := server.URL
	server.Close()

	setOpenAILLMEnv(t, closedURL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
}

func TestSimpleLLMCallMidResponseConnectionDrop(t *testing.T) {
	resetLLMClients(t)

	// 通过劫持底层连接并直接关闭，模拟对话过程中连接异常中断。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "no hijack", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}
		_ = conn.Close()
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
}

func TestSimpleLLMCallRateLimitedThenFails(t *testing.T) {
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"rate limit exceeded","type":"rate_limit_error"}}`))
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
	require.Equal(t, int32(fastRetry.attemptsPerModel), atomic.LoadInt32(&calls), "429 应当被重试")
}

func TestSimpleLLMCallAuthError(t *testing.T) {
	resetLLMClients(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid api key","type":"invalid_request_error"}}`))
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
}

func TestSimpleLLMCallContextLengthExceeded(t *testing.T) {
	resetLLMClients(t)

	var lastBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		lastBody = `{"error":{"message":"This model's maximum context length is 4096 tokens","type":"invalid_request_error","code":"context_length_exceeded"}}`
		_, _ = w.Write([]byte(lastBody))
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	_, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
	// 上游错误信息应被透传，便于运维定位 token 超限问题。
	require.Contains(t, err.Error(), "maximum context length")
}

func TestSimpleLLMCallTimeout(t *testing.T) {
	resetLLMClients(t)

	// 临时缩短调用超时，让 server 故意阻塞触发超时。
	oldTimeout := llmCallTimeout
	llmCallTimeout = 100 * time.Millisecond
	t.Cleanup(func() { llmCallTimeout = oldTimeout })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(time.Second):
		}
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	start := time.Now()
	_, err := simpleLLMCall("model-a", "hi", llmRetryConfig{
		attemptsPerModel: 1,
		delay:            time.Nanosecond,
		maxDelay:         time.Nanosecond,
	})
	elapsed := time.Since(start)

	require.Error(t, err)
	require.Contains(t, err.Error(), "all models failed")
	require.Less(t, elapsed, 2*time.Second, "超时应在 llmCallTimeout 后及时返回，而非等待整个上游耗时")
}

func TestSimpleLLMCallRecoversAfterTransientFailures(t *testing.T) {
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			http.Error(w, "temporary upstream failure", http.StatusBadGateway)
			return
		}
		writeChatCompletion(w, "recovered")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	result, err := simpleLLMCall("model-a", "hi", fastRetry)

	require.NoError(t, err)
	require.Equal(t, "recovered", result)
	require.Equal(t, int32(3), atomic.LoadInt32(&calls), "前两次失败后第三次应成功返回")
}

func TestSimpleLLMCallFailsOverToHealthyModel(t *testing.T) {
	resetLLMClients(t)

	var mu sync.Mutex
	requested := make([]string, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload chatCompletionRequest
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &payload)

		mu.Lock()
		requested = append(requested, payload.Model)
		mu.Unlock()

		// model-bad 始终失败，model-good 始终成功。
		if payload.Model == "model-bad" {
			http.Error(w, "bad model down", http.StatusServiceUnavailable)
			return
		}
		writeChatCompletion(w, "served by good model")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	result, err := simpleLLMCall("model-bad,model-good", "hi", fastRetry)

	require.NoError(t, err)
	require.Equal(t, "served by good model", result)

	mu.Lock()
	defer mu.Unlock()
	require.Contains(t, requested, "model-good", "应当在坏模型失败后切换到健康模型")
}

func TestCurrentLLMCallTimeoutDefaultAndEnv(t *testing.T) {
	old := llmCallTimeout
	llmCallTimeout = 0
	t.Cleanup(func() { llmCallTimeout = old })

	t.Setenv("FC_LLM_CALL_TIMEOUT", "")
	require.Equal(t, defaultLLMCallTimeout, currentLLMCallTimeout())
	require.Equal(t, 3*time.Minute, defaultLLMCallTimeout)

	t.Setenv("FC_LLM_CALL_TIMEOUT", "90")
	require.Equal(t, 90*time.Second, currentLLMCallTimeout())

	llmCallTimeout = 50 * time.Millisecond
	require.Equal(t, 50*time.Millisecond, currentLLMCallTimeout(), "测试注入的 llmCallTimeout 应覆盖环境变量")
}

func TestSimpleLLMCallContextCanceledDoesNotRetry(t *testing.T) {
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
		}
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := SimpleLLMCallContext(ctx, "model-a", "hi")
	elapsed := time.Since(start)

	require.Error(t, err)
	require.Less(t, elapsed, 2*time.Second)
	require.LessOrEqual(t, atomic.LoadInt32(&calls), int32(1), "请求取消后不应再重试占用并发额度")
}

func TestSimpleLLMCallContextCanceledDuringBackoffDoesNotRetry(t *testing.T) {
	resetLLMClients(t)

	firstCall := make(chan struct{})
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			close(firstCall)
		}
		http.Error(w, "retryable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)
	setOpenAILLMEnv(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := simpleLLMCallWithOptions(ctx, "model-a", "hi", llmRetryConfig{
			attemptsPerModel: 3,
			delay:            time.Second,
			maxDelay:         time.Second,
		}, nil)
		done <- err
	}()

	<-firstCall
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("cancellation did not interrupt retry backoff")
	}
	require.Equal(t, int32(1), atomic.LoadInt32(&calls))
}
