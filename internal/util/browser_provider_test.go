package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBrowserlessContentUsesBrowserEndpointForBrowserlessProvider(t *testing.T) {
	var captured BrowserRenderReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/content", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&captured))
		_, _ = w.Write([]byte("<html><body>rendered</body></html>"))
	}))
	defer server.Close()

	t.Setenv("FC_BROWSER_PROVIDER", "browserless-restful")
	t.Setenv("FC_BROWSER_ENDPOINT", server.URL)
	t.Setenv("FC_PUPPETEER_HTTP_ENDPOINT", "")

	content, err := GetBrowserlessContent("https://example.com/article", BrowserlessOptions{
		Timeout:   time.Second,
		WaitTime:  2 * time.Second,
		WaitUntil: "networkidle2",
	})

	require.NoError(t, err)
	assert.Equal(t, "<html><body>rendered</body></html>", content)
	assert.Equal(t, "https://example.com/article", captured.URL)
	assert.Equal(t, []string{"image"}, captured.RejectResourceTypes)
	assert.Equal(t, 2000, captured.WaitFor)
	require.NotNil(t, captured.GotoOptions)
	assert.Equal(t, "networkidle2", captured.GotoOptions.WaitUntil)
	assert.Equal(t, int64(1000), captured.GotoOptions.Timeout)
}

func TestResolveBrowserProviderConfigFallsBackToLegacyBrowserlessEndpoint(t *testing.T) {
	t.Setenv("FC_BROWSER_PROVIDER", "")
	t.Setenv("FC_BROWSER_ENDPOINT", "")
	t.Setenv("FC_PUPPETEER_HTTP_ENDPOINT", "http://legacy-browserless:3000")

	cfg := resolveBrowserProviderConfig(GetEnvClient())

	assert.Equal(t, BrowserProviderBrowserlessRESTful, cfg.Provider)
	assert.Equal(t, "http://legacy-browserless:3000", cfg.Endpoint)
}

func TestBrowserlessRestAliasIsNotSupported(t *testing.T) {
	assert.False(t, IsSupportedBrowserProvider("browserless-rest"))
}

func TestBuildEndpointURLPreservesCloakBrowserQueryParams(t *testing.T) {
	got, err := buildEndpointURL("http://service.cloakbrowser:9222?fingerprint=feedcraft&locale=zh-CN", "/json/version")

	require.NoError(t, err)
	assert.Equal(t, "http://service.cloakbrowser:9222/json/version?fingerprint=feedcraft&locale=zh-CN", got)
}

func TestGetBrowserlessContentWithCloakBrowserCDP(t *testing.T) {
	endpoint := os.Getenv("FC_TEST_CLOAK_BROWSER_ENDPOINT")
	if endpoint == "" {
		t.Skip("set FC_TEST_CLOAK_BROWSER_ENDPOINT to run the CloakBrowser CDP integration test")
	}

	t.Setenv("FC_BROWSER_PROVIDER", "cdp")
	t.Setenv("FC_BROWSER_ENDPOINT", endpoint)

	content, err := GetBrowserlessContent(
		"data:text/html,%3Chtml%3E%3Cbody%3E%3Ch1%3Ecloak-ok%3C%2Fh1%3E%3C%2Fbody%3E%3C%2Fhtml%3E",
		BrowserlessOptions{Timeout: 30 * time.Second, WaitUntil: "load"},
	)

	require.NoError(t, err)
	assert.True(t, strings.Contains(content, "cloak-ok"), content)
}

func TestResolveBrowserRenderTimeoutFromEnv(t *testing.T) {
	t.Setenv("FC_BROWSER_TIMEOUT", "")
	assert.Equal(t, DefaultBrowserRenderTimeout, ResolveBrowserRenderTimeout())

	t.Setenv("FC_BROWSER_TIMEOUT", "90s")
	assert.Equal(t, 90*time.Second, ResolveBrowserRenderTimeout())

	t.Setenv("FC_BROWSER_TIMEOUT", "120000")
	assert.Equal(t, 120*time.Second, ResolveBrowserRenderTimeout())
}

func TestResolveBrowserMaxConcurrencyFromEnv(t *testing.T) {
	t.Setenv("FC_BROWSER_MAX_CONCURRENCY", "")
	assert.Equal(t, DefaultBrowserMaxConcurrency, ResolveBrowserMaxConcurrency())

	t.Setenv("FC_BROWSER_MAX_CONCURRENCY", "4")
	assert.Equal(t, 4, ResolveBrowserMaxConcurrency())
}

func TestGetBrowserlessContentLimitsGlobalConcurrency(t *testing.T) {
	var inflight atomic.Int32
	var maxInflight atomic.Int32
	var started sync.WaitGroup
	started.Add(4)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := inflight.Add(1)
		for {
			prev := maxInflight.Load()
			if current <= prev || maxInflight.CompareAndSwap(prev, current) {
				break
			}
		}
		started.Done()
		time.Sleep(150 * time.Millisecond)
		inflight.Add(-1)
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	}))
	defer server.Close()

	t.Setenv("FC_BROWSER_PROVIDER", "browserless-restful")
	t.Setenv("FC_BROWSER_ENDPOINT", server.URL)
	t.Setenv("FC_BROWSER_MAX_CONCURRENCY", "2")
	resetBrowserRenderGateForTest()

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := GetBrowserlessContent(fmt.Sprintf("https://example.com/article-%d", i), BrowserlessOptions{
				Timeout: 2 * time.Second,
			})
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, maxInflight.Load(), int32(2))
	assert.Equal(t, int32(2), maxInflight.Load())
}
