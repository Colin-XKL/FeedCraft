package adapter

import (
	"FeedCraft/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// requireRedis 配置并探测 Redis，不可用时跳过依赖缓存的测试。
func requireRedis(t *testing.T) {
	t.Helper()
	t.Setenv("FC_REDIS_URI", "redis://127.0.0.1:6379")
	if err := util.TryCacheSetString("feedcraft_test_ping", "1", time.Minute); err != nil {
		t.Skipf("Redis 不可用，跳过缓存链路测试: %v", err)
	}
}

// uniquePrompt 生成进程内唯一前缀，避免命中上一轮测试遗留的 Redis 缓存。
func uniquePrompt(base string) string {
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano())
}

func TestCallLLMUsingContextCachesResult(t *testing.T) {
	requireRedis(t)
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		writeChatCompletion(w, "cached answer")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("summarize this")

	first, err := CallLLMUsingContext(prompt, "some article content", util.ContentProcessOption{})
	require.NoError(t, err)
	require.Equal(t, "cached answer", first)

	second, err := CallLLMUsingContext(prompt, "some article content", util.ContentProcessOption{})
	require.NoError(t, err)
	require.Equal(t, "cached answer", second)

	require.Equal(t, int32(1), atomic.LoadInt32(&calls), "相同 prompt+context 应命中缓存，仅请求上游一次")
}

func TestCallLLMUsingContextStripsBackticksFromContext(t *testing.T) {
	requireRedis(t)
	resetLLMClients(t)

	var captured chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, "ok")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("translate")
	_, err := CallLLMUsingContext(prompt, "danger `inline code` here", util.ContentProcessOption{})
	require.NoError(t, err)

	require.Len(t, captured.Messages, 1)
	content := captured.Messages[0].Content
	require.Contains(t, content, "danger inline code here", "context 内的反引号应被剥离")
	require.NotContains(t, content, "`inline code`", "context 内不应残留反引号包裹的内容")
}

func TestCallLLMUsingContextTemperatureSeparatesCache(t *testing.T) {
	requireRedis(t)
	resetLLMClients(t)

	var mu = make(chan struct{}, 1)
	mu <- struct{}{}
	temperatures := make([]float64, 0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload chatCompletionRequest
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &payload)
		<-mu
		temperatures = append(temperatures, payload.Temperature)
		mu <- struct{}{}
		writeChatCompletion(w, "ok")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("rate")
	context := "content for temperature test"

	zero := 0.0
	high := 0.9

	_, err := CallLLMUsingContext(prompt, context, util.ContentProcessOption{Temperature: &zero})
	require.NoError(t, err)
	_, err = CallLLMUsingContext(prompt, context, util.ContentProcessOption{Temperature: &high})
	require.NoError(t, err)

	<-mu
	captured := append([]float64(nil), temperatures...)
	mu <- struct{}{}

	require.Len(t, captured, 2, "不同 temperature 应使用不同缓存键，分别请求上游")
	require.Contains(t, captured, 0.0)
	require.Contains(t, captured, 0.9)
}

func TestCallLLMUsingContextErrorNotCached(t *testing.T) {
	requireRedis(t)
	resetLLMClients(t)

	// CallLLMUsingContext 内部走 SimpleLLMCall 使用默认重试配置，
	// 临时压缩退避时间，避免测试因真实退避而过慢。
	oldRetry := defaultLLMRetryConfig
	defaultLLMRetryConfig = fastRetry
	t.Cleanup(func() { defaultLLMRetryConfig = oldRetry })

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		http.Error(w, "upstream down", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("fail case")
	_, err := CallLLMUsingContext(prompt, "content", util.ContentProcessOption{})
	require.Error(t, err)

	before := atomic.LoadInt32(&calls)
	_, err = CallLLMUsingContext(prompt, "content", util.ContentProcessOption{})
	require.Error(t, err)
	require.Greater(t, atomic.LoadInt32(&calls), before, "失败结果不应被缓存，二次调用应再次请求上游")
}
