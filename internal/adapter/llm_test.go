package adapter

import (
	"FeedCraft/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
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

func setupAdapterMiniredis(t *testing.T) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Setenv("FC_REDIS_URI", fmt.Sprintf("redis://%s", mr.Addr()))
	t.Cleanup(mr.Close)
}

func TestCallLLMUsingContextCompactsAndRestoresBase64AndLongURLs(t *testing.T) {
	setupAdapterMiniredis(t)
	resetLLMClients(t)

	largeBase64 := "data:image/png;base64," + strings.Repeat("QUJDREVGR0hJSktMTU5PUA==", 500)
	longURL := "https://cdn.example.com/articles/view?" + strings.Repeat("tracking=1234567890&", 10) + "ref=rss"

	var captured chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		// Echo back a translated version that keeps the placeholders intact.
		writeChatCompletion(w, `<p>已翻译</p><img src="https://fc.local/m/1"><a href="https://fc.local/u/1">链接</a>`)
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("translate")
	inputHTML := fmt.Sprintf(`<p>Original</p><img src="%s"><a href="%s">Link</a>`, largeBase64, longURL)

	result, err := CallLLMUsingContext(prompt, inputHTML, util.ContentProcessOption{})
	require.NoError(t, err)

	require.Len(t, captured.Messages, 1)
	sentPrompt := captured.Messages[0].Content
	require.NotContains(t, sentPrompt, "base64", "base64 数据不应透传给上游 LLM")
	require.NotContains(t, sentPrompt, longURL, "过长 URL 不应透传给上游 LLM")
	require.Contains(t, sentPrompt, "https://fc.local/m/1")
	require.Contains(t, sentPrompt, "https://fc.local/u/1")

	require.Contains(t, result, `<p>已翻译</p>`)
	require.Contains(t, result, `src="`+largeBase64+`"`, "返回结果应原样还原 base64 媒体资源")
	require.Contains(t, result, `href="`+longURL+`"`, "返回结果应原样还原超长 URL")
}

func TestCallLLMUsingContextCacheRestoresDistinctPayloads(t *testing.T) {
	setupAdapterMiniredis(t)
	resetLLMClients(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		writeChatCompletion(w, `<p>译文</p><img src="https://fc.local/m/1">`)
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("translate-cache-distinct")
	imgA := "data:image/png;base64," + strings.Repeat("AAAA", 100)
	imgB := "data:image/png;base64," + strings.Repeat("BBBB", 100)

	resA, err := CallLLMUsingContext(prompt, `<p>Same text</p><img src="`+imgA+`">`, util.ContentProcessOption{})
	require.NoError(t, err)
	resB, err := CallLLMUsingContext(prompt, `<p>Same text</p><img src="`+imgB+`">`, util.ContentProcessOption{})
	require.NoError(t, err)

	require.Equal(t, int32(1), atomic.LoadInt32(&calls), "正文一致且占位符结构一致时可复用缓存，仅请求上游一次")
	require.Contains(t, resA, imgA)
	require.NotContains(t, resA, imgB)
	require.Contains(t, resB, imgB)
	require.NotContains(t, resB, imgA)
}

func TestCallLLMUsingContextRemoveImageDropsBase64MarkdownAndSourceBeforeCompaction(t *testing.T) {
	setupAdapterMiniredis(t)
	resetLLMClients(t)

	imgData := "data:image/png;base64," + strings.Repeat("QUJD", 50)
	var captured chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, "true")
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("check-remove-image")
	context := fmt.Sprintf(
		`<p>Body</p><picture><source srcset="%s"></picture><img src=%s> ![alt](%s)`,
		imgData,
		imgData,
		imgData,
	)

	res, err := CallLLMUsingContext(prompt, context, util.ContentProcessOption{RemoveImage: true})
	require.NoError(t, err)
	require.Equal(t, "true", res)
	require.Len(t, captured.Messages, 1)
	sentPrompt := captured.Messages[0].Content
	require.NotContains(t, sentPrompt, "base64")
	require.NotContains(t, sentPrompt, "https://fc.local/m/", "被 RemoveImage 移除的图片不应以占位符形式发送给 LLM")
}

func TestCallLLMUsingContextReservesPromptPlaceholders(t *testing.T) {
	setupAdapterMiniredis(t)
	resetLLMClients(t)

	imgData := "data:image/png;base64," + strings.Repeat("QUJD", 50)
	var captured chatCompletionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &captured))
		writeChatCompletion(w, `Keep https://fc.local/m/1 and <img src="https://fc.local/m/2">`)
	}))
	t.Cleanup(server.Close)

	setOpenAILLMEnv(t, server.URL)
	t.Setenv("FC_LLM_API_MODEL", "model-a")

	prompt := uniquePrompt("Keep literal https://fc.local/m/1 in output")
	res, err := CallLLMUsingContext(prompt, `<img src="`+imgData+`">`, util.ContentProcessOption{})
	require.NoError(t, err)

	require.Len(t, captured.Messages, 1)
	require.Contains(t, captured.Messages[0].Content, `<img src="https://fc.local/m/2">`)
	require.Contains(t, res, "Keep https://fc.local/m/1")
	require.Contains(t, res, `<img src="`+imgData+`">`)
}
