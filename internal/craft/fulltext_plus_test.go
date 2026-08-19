package craft

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFulltextPlusConfigDefaults(t *testing.T) {
	config := parseFulltextPlusConfig(map[string]string{})
	assert.Equal(t, 0, config.Wait)
	assert.Equal(t, "load", config.Mode)
}

func TestBuildFulltextPlusOptionsUsesBrowserTimeout(t *testing.T) {
	opts := buildFulltextPlusBrowserOptions(FulltextPlusConfig{Mode: "load"})
	assert.Equal(t, "load", opts.WaitUntil)
	assert.Equal(t, util.DefaultBrowserRenderTimeout, opts.Timeout)
	assert.Equal(t, time.Duration(0), opts.WaitTime)

	opts = buildFulltextPlusBrowserOptions(FulltextPlusConfig{Wait: 90, Mode: "networkidle2"})
	assert.Equal(t, 90*time.Second, opts.WaitTime)
	assert.Equal(t, 100*time.Second, opts.Timeout)
	assert.Equal(t, "networkidle2", opts.WaitUntil)
}

func TestFulltextPlusProcessor_FallsBackToHTTPWhenBrowserFails(t *testing.T) {
	setupTestRedis(t)

	articleHTML := `<!DOCTYPE html><html><head><title>Fallback Article</title></head><body>
<article><h1>Fallback Article</h1>
<p>This paragraph is long enough for readability to treat this page as a real article with substantial content for extraction.</p>
<p>Another paragraph keeps the extractor confident that the main body lives here instead of in navigation chrome.</p>
</article></body></html>`
	articleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(articleHTML))
	}))
	t.Cleanup(articleServer.Close)

	original := fulltextPlusExtractFunc
	fulltextPlusExtractFunc = func(url string, options util.BrowserlessOptions) (string, error) {
		return "", fmt.Errorf("Post \"http://service.browserless:3000/content\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)")
	}
	t.Cleanup(func() { fulltextPlusExtractFunc = original })

	processor := newFulltextPlusProcessor("https://example.com/feed.xml", FulltextPlusConfig{Mode: "load"})
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "fallback-" + t.Name(), Link: articleServer.URL},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	require.NotEmpty(t, result.Articles[0].Content)
	assert.Contains(t, result.Articles[0].Content, "substantial content for extraction")
}

func TestFulltextPlusProcessor_OpensCircuitAfterConsecutiveBrowserFailures(t *testing.T) {
	setupTestRedis(t)

	articleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusBadGateway)
	}))
	t.Cleanup(articleServer.Close)

	var calls atomic.Int32
	original := fulltextPlusExtractFunc
	fulltextPlusExtractFunc = func(url string, options util.BrowserlessOptions) (string, error) {
		calls.Add(1)
		return "", fmt.Errorf("browser timeout")
	}
	t.Cleanup(func() { fulltextPlusExtractFunc = original })

	processor := newFulltextPlusProcessor("https://example.com/feed.xml", FulltextPlusConfig{Mode: "load"})
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "c1-" + t.Name(), Link: articleServer.URL + "/1"},
			{Title: "c2-" + t.Name(), Link: articleServer.URL + "/2"},
			{Title: "c3-" + t.Name(), Link: articleServer.URL + "/3"},
			{Title: "c4-" + t.Name(), Link: articleServer.URL + "/4"},
		},
	}

	_, _ = processor.Process(context.Background(), feed)
	assert.Equal(t, int32(browserFailCircuitThreshold), calls.Load())
}
