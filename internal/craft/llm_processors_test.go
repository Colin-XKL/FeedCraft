package craft

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"FeedCraft/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testTinyBase64PNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

// TestGetArticleContentForPrompt_StripsBase64Images verifies that the LLM
// preparation path (used by summary, introduction, etc.) strips inline base64
// images before converting HTML to Markdown for the LLM context.
func TestGetArticleContentForPrompt_StripsBase64Images(t *testing.T) {
	article := &model.CraftArticle{
		Title: "Test Article",
		Link:  "https://example.com/post/1",
	}
	original := `<h1>Weekly Roundup</h1>` +
		`<p>Read on for details.</p>` +
		`<img src="data:image/png;base64,` + testTinyBase64PNG + `" alt="cover">` +
		`<p>End of article.</p>`

	result := getArticleContentForPrompt(article, original)

	assert.NotContains(t, result, "base64", "base64 image data must not be sent to the LLM")
	assert.NotContains(t, result, "data:image", "base64 data URI must be stripped")
	assert.Contains(t, result, "Weekly Roundup")
	assert.Contains(t, result, "Read on for details.")
	assert.Contains(t, result, "End of article.")
	// Stripped output should be far smaller than the raw base64 blob.
	assert.Less(t, len(result), len(testTinyBase64PNG))
}

// TestGetArticleContentForPrompt_KeepsNormalImages verifies that regular
// (non-base64) image references are preserved after the cleanup.
func TestGetArticleContentForPrompt_KeepsNormalImages(t *testing.T) {
	article := &model.CraftArticle{
		Title: "Test Article",
		Link:  "https://example.com/post/2",
	}
	original := `<p>See diagram below.</p><img src="https://example.com/pic.png" alt="diagram">`

	result := getArticleContentForPrompt(article, original)

	assert.True(t,
		strings.Contains(result, "example.com/pic.png") || strings.Contains(result, "diagram"),
		"normal image reference should be preserved",
	)
}

type processorResult struct {
	feed *model.CraftFeed
	err  error
}

func observeConcurrentStart(started <-chan struct{}, release chan struct{}) bool {
	<-started
	select {
	case <-started:
		close(release)
		return true
	case <-time.After(200 * time.Millisecond):
		close(release)
		return false
	}
}

func TestArticleTextTransformProcessor_Process_RunsConcurrentlyAndPreservesOrder(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "1")

	started := make(chan struct{}, 3)
	release := make(chan struct{})
	processor := &ArticleTextTransformProcessor{
		CraftName: "concurrent transform",
		Mutate: func(_ context.Context, article *model.CraftArticle) error {
			started <- struct{}{}
			<-release
			article.Description = "processed:" + article.Title
			return nil
		},
	}
	feed := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "first"},
		{Title: "second"},
		{Title: "third"},
	}}

	done := make(chan processorResult, 1)
	go func() {
		processed, err := processor.Process(context.Background(), feed)
		done <- processorResult{feed: processed, err: err}
	}()

	concurrent := observeConcurrentStart(started, release)
	result := <-done

	require.NoError(t, result.err)
	require.Len(t, result.feed.Articles, 3)
	assert.True(t, concurrent, "a single feed should start more than one article concurrently")
	assert.Equal(t, "processed:first", result.feed.Articles[0].Description)
	assert.Equal(t, "processed:second", result.feed.Articles[1].Description)
	assert.Equal(t, "processed:third", result.feed.Articles[2].Description)
}

func TestArticleTextTransformProcessor_Process_PartialFailureKeepsSuccessfulResults(t *testing.T) {
	processor := &ArticleTextTransformProcessor{
		CraftName: "partial failure",
		Mutate: func(_ context.Context, article *model.CraftArticle) error {
			if article.Title == "failed" {
				return errors.New("LLM unavailable")
			}
			article.Description = "processed"
			return nil
		},
	}
	feed := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "successful"},
		{Title: "failed"},
		nil,
	}}

	processed, err := processor.Process(context.Background(), feed)

	require.NoError(t, err)
	require.Len(t, processed.Articles, 3)
	assert.Equal(t, "processed", processed.Articles[0].Description)
	assert.Empty(t, processed.Articles[1].Description)
	assert.Nil(t, processed.Articles[2])
}

func TestArticleTextTransformProcessor_Process_AllFailuresReturnError(t *testing.T) {
	processor := &ArticleTextTransformProcessor{
		CraftName: "total failure",
		Mutate: func(_ context.Context, _ *model.CraftArticle) error {
			return errors.New("LLM unavailable")
		},
	}
	feed := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "first"},
		{Title: "second"},
	}}

	processed, err := processor.Process(context.Background(), feed)

	require.Error(t, err)
	assert.Nil(t, processed)
	assert.Contains(t, err.Error(), "all items failed to process")
}

func TestArticlePredicateProcessor_Process_RunsConcurrentlyAndPreservesOrder(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "1")

	started := make(chan struct{}, 4)
	release := make(chan struct{})
	processor := &ArticlePredicateProcessor{
		CraftName: "concurrent predicate",
		Match: func(_ context.Context, article *model.CraftArticle) (bool, error) {
			started <- struct{}{}
			<-release
			return article.Title == "drop", nil
		},
	}
	feed := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "keep-first"},
		{Title: "drop"},
		{Title: "keep-last"},
	}}

	done := make(chan processorResult, 1)
	go func() {
		processed, err := processor.Process(context.Background(), feed)
		done <- processorResult{feed: processed, err: err}
	}()

	concurrent := observeConcurrentStart(started, release)
	result := <-done

	require.NoError(t, result.err)
	require.Len(t, result.feed.Articles, 2)
	assert.True(t, concurrent, "a single feed should evaluate more than one article concurrently")
	assert.Equal(t, "keep-first", result.feed.Articles[0].Title)
	assert.Equal(t, "keep-last", result.feed.Articles[1].Title)
}

func TestArticlePredicateProcessor_Process_FailureKeepsArticleAndNilIsSkipped(t *testing.T) {
	processor := &ArticlePredicateProcessor{
		CraftName: "predicate failure",
		Match: func(_ context.Context, article *model.CraftArticle) (bool, error) {
			switch article.Title {
			case "failed":
				return false, errors.New("LLM unavailable")
			case "drop":
				return true, nil
			default:
				return false, nil
			}
		},
	}
	feed := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "keep"},
		nil,
		{Title: "failed"},
		{Title: "drop"},
	}}

	processed, err := processor.Process(context.Background(), feed)

	require.NoError(t, err)
	require.Len(t, processed.Articles, 2)
	assert.Equal(t, "keep", processed.Articles[0].Title)
	assert.Equal(t, "failed", processed.Articles[1].Title)
}

func TestBeautifyArticleContent_CompactsAndRestoresBase64AndLongURLs(t *testing.T) {
	origCaller := llmCaller
	t.Cleanup(func() { llmCaller = origCaller })

	largeBase64 := "data:image/png;base64," + strings.Repeat("QUJDREVGR0hJSktMTU5PUA==", 500)
	longURL := "https://cdn.example.com/articles/view?" + strings.Repeat("tracking=1234567890&", 10) + "ref=rss"

	var capturedPrompt string
	llmCaller = func(model string, promptInput string) (string, error) {
		capturedPrompt = promptInput
		return "Clean paragraph.\n\n![cover](https://fc.local/m/1)\n\n[Read more](https://fc.local/u/1)", nil
	}

	inputHTML := `<p>Messy   paragraph.</p><img src="` + largeBase64 + `" alt="cover"><a href="` + longURL + `">Read more</a>`
	outHTML, err := beautifyArticleContent(inputHTML, "Beautify this article")
	require.NoError(t, err)

	assert.NotContains(t, capturedPrompt, "base64", "base64 media payload must not be sent to LLM")
	assert.NotContains(t, capturedPrompt, longURL, "overly long URL must not be sent to LLM")
	assert.Contains(t, capturedPrompt, "https://fc.local/m/1")
	assert.Contains(t, capturedPrompt, "https://fc.local/u/1")

	assert.Contains(t, outHTML, `src="`+largeBase64+`"`, "restored HTML must contain original base64 image")
	assert.Contains(t, outHTML, `href="`+longURL+`"`, "restored HTML must contain original long URL")
	assert.NotContains(t, outHTML, "https://fc.local/")
}

func TestBeautifyContentProcessor_CachesCompactRepresentationWithoutBase64(t *testing.T) {
	redisSrv := setupTestRedis(t)
	origCaller := llmCaller
	t.Cleanup(func() { llmCaller = origCaller })

	var calls int
	llmCaller = func(model string, promptInput string) (string, error) {
		calls++
		return "Beautified.\n\n![img](https://fc.local/m/1)", nil
	}

	imgA := "data:image/png;base64," + strings.Repeat("AAAA", 200)
	imgB := "data:image/png;base64," + strings.Repeat("BBBB", 200)

	processor := newBeautifyContentProcessor("beautify cache compact test")
	feedA := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "Same Title", Content: `<p>Same body</p><img src="` + imgA + `">`},
	}}
	feedB := &model.CraftFeed{Articles: []*model.CraftArticle{
		{Title: "Same Title", Content: `<p>Same body</p><img src="` + imgB + `">`},
	}}

	outA, err := processor.Process(context.Background(), feedA)
	require.NoError(t, err)
	outB, err := processor.Process(context.Background(), feedB)
	require.NoError(t, err)

	assert.Equal(t, 1, calls, "second article with same compacted structure should hit craft cache")
	assert.Contains(t, outA.Articles[0].Content, imgA)
	assert.NotContains(t, outA.Articles[0].Content, imgB)
	assert.Contains(t, outB.Articles[0].Content, imgB)
	assert.NotContains(t, outB.Articles[0].Content, imgA)

	// Verify Redis cache entry does not store raw base64 payload.
	keys := redisSrv.server.Keys()
	require.NotEmpty(t, keys)
	for _, k := range keys {
		val, getErr := redisSrv.server.Get(k)
		require.NoError(t, getErr)
		assert.NotContains(t, val, "base64", "craft cache must store compact HTML, not raw base64 media")
		assert.Contains(t, val, "https://fc.local/m/1")
	}
}

func TestBeautifyArticleContent_DoesNotCompactActiveDataHTMLLinks(t *testing.T) {
	origCaller := llmCaller
	t.Cleanup(func() { llmCaller = origCaller })

	activeDataHTML := "data:text/html;base64," + strings.Repeat("PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==", 8)
	var capturedPrompt string
	llmCaller = func(model string, promptInput string) (string, error) {
		capturedPrompt = promptInput
		return "[click](https://fc.local/m/1)", nil
	}

	outHTML, err := beautifyArticleContent(`<a href="`+activeDataHTML+`">click</a>`, "Beautify")
	require.NoError(t, err)
	assert.NotContains(t, capturedPrompt, "https://fc.local/m/", "active data:text/html URI must not be rewritten to placeholder")
	assert.NotContains(t, outHTML, "data:text/html", "Restore must not inject active data:text/html URI into rendered HTML")
}
