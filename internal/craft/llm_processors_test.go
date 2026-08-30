package craft

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
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

func TestArticleTextTransformProcessor_Process_RunsArticlesConcurrently(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "4")

	var inFlight atomic.Int32
	var maxInFlight atomic.Int32
	proc := &ArticleTextTransformProcessor{
		CraftName: "test-concurrent-transform",
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			n := inFlight.Add(1)
			defer inFlight.Add(-1)
			for {
				old := maxInFlight.Load()
				if n <= old || maxInFlight.CompareAndSwap(old, n) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			article.Description = "ok:" + article.Title
			return nil
		},
	}

	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "a"}, {Title: "b"}, {Title: "c"}, {Title: "d"},
		},
	}

	out, err := proc.Process(context.Background(), feed)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.Articles, 4)
	assert.Equal(t, "ok:a", out.Articles[0].Description)
	assert.Equal(t, "ok:b", out.Articles[1].Description)
	assert.Equal(t, "ok:c", out.Articles[2].Description)
	assert.Equal(t, "ok:d", out.Articles[3].Description)
	assert.GreaterOrEqual(t, maxInFlight.Load(), int32(2), "articles in a single feed should be processed concurrently")
}

func TestArticleTextTransformProcessor_Process_PartialFailureKeepsSuccesses(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "4")

	proc := &ArticleTextTransformProcessor{
		CraftName: "test-partial-fail",
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			if article.Title == "b" {
				return errors.New("llm failed")
			}
			article.Description = "ok"
			return nil
		},
	}
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "a"}, {Title: "b"}, {Title: "c"},
		},
	}

	out, err := proc.Process(context.Background(), feed)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "ok", out.Articles[0].Description)
	assert.Equal(t, "", out.Articles[1].Description)
	assert.Equal(t, "ok", out.Articles[2].Description)
}

func TestArticleTextTransformProcessor_Process_AllFailuresReturnError(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "3")

	proc := &ArticleTextTransformProcessor{
		CraftName: "test-all-fail",
		Mutate: func(ctx context.Context, article *model.CraftArticle) error {
			return errors.New("llm failed")
		},
	}
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "a"}, {Title: "b"},
		},
	}

	out, err := proc.Process(context.Background(), feed)
	require.Error(t, err)
	assert.Nil(t, out)
	assert.Contains(t, err.Error(), "all items failed")
}

func TestArticlePredicateProcessor_Process_RunsArticlesConcurrentlyAndPreservesOrder(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "4")

	var inFlight atomic.Int32
	var maxInFlight atomic.Int32
	proc := &ArticlePredicateProcessor{
		CraftName: "test-concurrent-predicate",
		Match: func(ctx context.Context, article *model.CraftArticle) (bool, error) {
			n := inFlight.Add(1)
			defer inFlight.Add(-1)
			for {
				old := maxInFlight.Load()
				if n <= old || maxInFlight.CompareAndSwap(old, n) {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
			if article.Title == "drop" {
				return true, nil
			}
			if article.Title == "err" {
				return false, errors.New("llm failed")
			}
			return false, nil
		},
	}

	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Title: "keep-1"},
			{Title: "drop"},
			{Title: "err"},
			{Title: "keep-2"},
		},
	}

	out, err := proc.Process(context.Background(), feed)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.Articles, 3)
	assert.Equal(t, "keep-1", out.Articles[0].Title)
	assert.Equal(t, "err", out.Articles[1].Title)
	assert.Equal(t, "keep-2", out.Articles[2].Title)
	assert.GreaterOrEqual(t, maxInFlight.Load(), int32(2), "predicate evaluation should run concurrently")
}

func TestForEachArticleConcurrently_RespectsLLMMaxConcurrency(t *testing.T) {
	t.Setenv("FC_LLM_MAX_CONCURRENCY", "2")

	var inFlight atomic.Int32
	var maxInFlight atomic.Int32
	articles := []*model.CraftArticle{
		{Title: "a"}, {Title: "b"}, {Title: "c"}, {Title: "d"},
	}

	forEachArticleConcurrently(articles, func(_ int, article *model.CraftArticle) {
		n := inFlight.Add(1)
		defer inFlight.Add(-1)
		for {
			old := maxInFlight.Load()
			if n <= old || maxInFlight.CompareAndSwap(old, n) {
				break
			}
		}
		time.Sleep(40 * time.Millisecond)
	})

	assert.Equal(t, int32(2), maxInFlight.Load())
}

func TestSanitizeTranslatedTitle_StripsLeakedLabels(t *testing.T) {
	assert.Equal(t, "香港楼市降温", sanitizeTranslatedTitle("文章内容：香港楼市降温"))
	assert.Equal(t, "Hong Kong housing cools", sanitizeTranslatedTitle("Article Content: Hong Kong housing cools"))
	assert.Equal(t, "只留第一行", sanitizeTranslatedTitle("标题：只留第一行\n第二行"))
	assert.Equal(t, "干净标题", sanitizeTranslatedTitle("「干净标题」"))
}
