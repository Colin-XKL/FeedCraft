package craft

import (
	"context"
	"fmt"
	"testing"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestResolveCraftAtoms_FlowAndCustomAtom(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	require.NoError(t, dao.CreateCraftAtom(db, &dao.CraftAtom{
		Name:         "limit-five",
		TemplateName: "limit",
		Params:       map[string]string{"num": "5"},
	}))
	require.NoError(t, dao.CreateCraftFlow(db, &dao.CraftFlow{
		Name: "flow-a",
		CraftFlowConfig: []dao.CraftFlowItem{
			{CraftName: "limit-five"},
			{CraftName: "guid-fix"},
		},
	}))

	atoms, err := ResolveCraftAtoms(db, "flow-a")
	require.NoError(t, err)
	require.Len(t, atoms, 2)
	assert.Equal(t, "limit", atoms[0].TemplateName)
	assert.Equal(t, "guid-fix", atoms[1].TemplateName)
	assert.Equal(t, "5", atoms[0].Params["num"])
}

func TestBuildProcessor_UsesNativeProcessors(t *testing.T) {
	db := newCraftRuntimeTestDB(t)
	processor, err := BuildProcessor(db, "limit,time-limit,guid-fix,relative-link-fix,cleanup,fulltext,fulltext-plus", "https://example.com/feed.xml")
	require.NoError(t, err)
	require.NotNil(t, processor)

	flow, ok := processor.(*engine.FlowCraftProcessor)
	require.True(t, ok)
	require.Len(t, flow.Processors, 7)
	assert.IsType(t, &LimitProcessor{}, flow.Processors[0])
	assert.IsType(t, &TimeLimitProcessor{}, flow.Processors[1])
	assert.IsType(t, &GUIDFixProcessor{}, flow.Processors[2])
	assert.IsType(t, &RelativeLinkFixProcessor{}, flow.Processors[3])
	assert.IsType(t, &CleanupProcessor{}, flow.Processors[4])
	assert.IsType(t, &FulltextProcessor{}, flow.Processors[5])
	assert.IsType(t, &FulltextPlusProcessor{}, flow.Processors[6])
}

func TestNativeProcessors_EndToEnd(t *testing.T) {
	now := time.Now()
	processor := &engine.FlowCraftProcessor{
		Processors: []engine.FeedProcessor{
			&KeywordProcessor{
				Mode:       KeywordIncludeMode,
				MatchScope: KeywordMatchAll,
				Keywords:   []string{"keep"},
			},
			&GUIDFixProcessor{},
			&RelativeLinkFixProcessor{OriginalFeedURL: "https://example.com/feed.xml"},
			&LimitProcessor{MaxItems: 1},
			&TimeLimitProcessor{
				Days: 7,
				Now:  func() time.Time { return now },
			},
		},
	}

	feed := &model.CraftFeed{
		Link: "https://example.com",
		Articles: []*model.CraftArticle{
			{
				Title:       "keep article",
				Content:     "keep this",
				Description: "keep this",
				Link:        "/article-1",
				Created:     now,
			},
			{
				Title:       "drop article",
				Content:     "drop this",
				Description: "drop this",
				Link:        "/article-2",
				Created:     now,
			},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	require.Len(t, result.Articles, 1)
	assert.Equal(t, "https://example.com/article-1", result.Articles[0].Link)
	assert.NotEmpty(t, result.Articles[0].Id)
}

func TestCleanupProcessor_UsesDescriptionFallback(t *testing.T) {
	original := cleanupTransformFunc
	cleanupTransformFunc = func(content string, domain string) (string, error) {
		return fmt.Sprintf("%s|%s", domain, content), nil
	}
	t.Cleanup(func() { cleanupTransformFunc = original })

	processor := newCleanupProcessor()
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{
				Title:       "article",
				Link:        "https://example.com/post",
				Description: "<p>fallback</p>",
			},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	require.Len(t, result.Articles, 1)
	assert.Equal(t, "https://example.com|<p>fallback</p>", result.Articles[0].Content)
	assert.Equal(t, "https://example.com|<p>fallback</p>", result.Articles[0].Description)
	assert.Empty(t, feed.Articles[0].Content)
}

func TestFulltextProcessor_PartialFailureAndRelativeLinkFix(t *testing.T) {
	original := fulltextExtractFunc
	fulltextExtractFunc = func(url string, timeout time.Duration) (string, error) {
		if url == "https://example.com/fail" {
			return "", fmt.Errorf("boom")
		}
		return "fulltext:" + url, nil
	}
	t.Cleanup(func() { fulltextExtractFunc = original })

	processor := newFulltextProcessor("https://example.com/feed.xml")
	feed := &model.CraftFeed{
		Link: "https://example.com",
		Articles: []*model.CraftArticle{
			{Title: "ok", Link: "/ok"},
			{Title: "bad", Link: "/fail"},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	require.Len(t, result.Articles, 2)
	assert.Equal(t, "https://example.com/ok", result.Articles[0].Link)
	assert.Equal(t, "fulltext:https://example.com/ok", result.Articles[0].Content)
	assert.Equal(t, "fulltext:https://example.com/ok", result.Articles[0].Description)
	assert.Equal(t, "https://example.com/fail", result.Articles[1].Link)
	assert.Empty(t, result.Articles[1].Content)
}

func TestFulltextProcessor_AllFailureReturnsError(t *testing.T) {
	original := fulltextExtractFunc
	fulltextExtractFunc = func(url string, timeout time.Duration) (string, error) {
		return "", fmt.Errorf("always fail")
	}
	t.Cleanup(func() { fulltextExtractFunc = original })

	processor := newFulltextProcessor("https://example.com/feed.xml")
	feed := &model.CraftFeed{
		Link: "https://example.com",
		Articles: []*model.CraftArticle{
			{Title: "a", Link: "/a"},
			{Title: "b", Link: "/b"},
		},
	}

	_, err := processor.Process(context.Background(), feed)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all items failed to process")
}

func TestFulltextPlusProcessor_UsesConfiguredOptions(t *testing.T) {
	original := fulltextPlusExtractFunc
	var capturedURL string
	var capturedOptions util.BrowserlessOptions
	fulltextPlusExtractFunc = func(url string, options util.BrowserlessOptions) (string, error) {
		capturedURL = url
		capturedOptions = options
		return "rendered", nil
	}
	t.Cleanup(func() { fulltextPlusExtractFunc = original })

	processor := newFulltextPlusProcessor("https://example.com/feed.xml", FulltextPlusConfig{
		Wait: 42,
		Mode: "networkidle0",
	})
	feed := &model.CraftFeed{
		Link: "https://example.com",
		Articles: []*model.CraftArticle{
			{Title: "a", Link: "/article"},
		},
	}

	result, err := processor.Process(context.Background(), feed)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/article", capturedURL)
	assert.Equal(t, "networkidle0", capturedOptions.WaitUntil)
	assert.Equal(t, 42*time.Second, capturedOptions.WaitTime)
	assert.Equal(t, 52*time.Second, capturedOptions.Timeout)
	assert.Equal(t, "rendered", result.Articles[0].Content)
}

func newCraftRuntimeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&dao.CraftAtom{}, &dao.CraftFlow{}))
	return db
}
