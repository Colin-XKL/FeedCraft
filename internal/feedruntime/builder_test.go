package feedruntime

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/engine"
	"FeedCraft/internal/model"
	"FeedCraft/internal/source"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBuildProviderFromInput_RecipeURI(t *testing.T) {
	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.CustomRecipeV2{
		ID:           "recipe-1",
		Craft:        "proxy",
		SourceType:   string(constant.SourceRSS),
		SourceConfig: `{"type":"rss","http_fetcher":{"url":"https://example.com/feed.xml"}}`,
	}).Error)

	builder := NewBuilder(db)
	provider, err := builder.BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindURI,
		URI:  "feedcraft://recipe/recipe-1",
	}, nil)
	require.NoError(t, err)
	assert.IsType(t, &RecipeProvider{}, provider)
}

func TestBuildRecipeProvider(t *testing.T) {
	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.CustomRecipeV2{
		ID:           "recipe-1",
		Craft:        "proxy",
		SourceType:   string(constant.SourceRSS),
		SourceConfig: `{"type":"rss","http_fetcher":{"url":"https://example.com/feed.xml"}}`,
	}).Error)

	builder := NewBuilder(db)
	provider, err := builder.BuildRecipeProvider(context.Background(), "recipe-1")
	require.NoError(t, err)
	assert.IsType(t, &RecipeProvider{}, provider)
}

func TestBuildProviderFromInput_HTTPURL(t *testing.T) {
	builder := NewBuilder(newTestDB(t))
	provider, err := builder.BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindURI,
		URI:  "https://example.com/feed.xml",
	}, nil)
	require.NoError(t, err)
	rawProvider, ok := provider.(*RawFeedProvider)
	require.True(t, ok)
	assert.Equal(t, "https://example.com/feed.xml", rawProvider.URL)
}

func TestBuildProviderFromInput_InboxURIUsesBuilderDB(t *testing.T) {
	db := newTestDB(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:          "inbox-1",
		Title:       "Inbox Feed",
		Description: "Inbox description",
		MaxItems:    100,
	}))
	require.NoError(t, db.Create(&dao.InboxItem{
		InboxID:     "inbox-1",
		ItemID:      "item-1",
		Title:       "Inbox Item",
		URL:         "https://example.com/inbox-item",
		Content:     "<p>Inbox content</p>",
		Summary:     "Inbox summary",
		PublishedAt: time.Unix(1700000000, 0),
		CreatedAt:   time.Unix(1700000000, 0),
	}).Error)

	provider, err := NewBuilder(db).BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindURI,
		URI:  "feedcraft://inbox/inbox-1",
	}, nil)
	require.NoError(t, err)
	assert.IsType(t, &InboxProvider{}, provider)

	feed, err := provider.Fetch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "Inbox Feed", feed.Title)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "Inbox Item", feed.Articles[0].Title)
}

func TestBuildProviderFromInput_SourceConfig(t *testing.T) {
	const testSourceType = constant.SourceType("unit_test_source")
	registerTestSource(t, testSourceType, func(cfg *config.SourceConfig) (source.Source, error) {
		return &stubSource{baseURL: cfg.HttpFetcher.URL}, nil
	})

	builder := NewBuilder(newTestDB(t))
	provider, err := builder.BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindSource,
		SourceConfig: &config.SourceConfig{
			Type: testSourceType,
			HttpFetcher: &config.HttpFetcherConfig{
				URL: "https://example.com/source",
			},
		},
	}, nil)
	require.NoError(t, err)

	feed, err := provider.Fetch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "stub-feed", feed.Title)
}

func TestBuildRecipe_AppliesCraftProcessor(t *testing.T) {
	const testSourceType = constant.SourceType("unit_test_recipe_source")
	registerTestSource(t, testSourceType, func(cfg *config.SourceConfig) (source.Source, error) {
		return &stubSource{
			baseURL:          cfg.HttpFetcher.URL,
			itemLinkOverride: "/relative-item",
		}, nil
	})

	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.CustomRecipeV2{
		ID:         "recipe-relative-link-fix",
		Craft:      "relative-link-fix",
		SourceType: string(testSourceType),
		SourceConfig: `{
			"type":"unit_test_recipe_source",
			"http_fetcher":{"url":"https://example.com/base/feed.xml"}
		}`,
	}).Error)

	builder := NewBuilder(db)
	provider, err := builder.BuildRecipeProvider(context.Background(), "recipe-relative-link-fix")
	require.NoError(t, err)

	feed, err := provider.Fetch(context.Background())
	require.NoError(t, err)
	require.Len(t, feed.Articles, 1)
	assert.Equal(t, "https://example.com/relative-item", feed.Articles[0].Link)
}

func TestBuildRecipe_UsesSourceInputSpecCompatibility(t *testing.T) {
	const testSourceType = constant.SourceType("unit_test_recipe_source_compat")
	registerTestSource(t, testSourceType, func(cfg *config.SourceConfig) (source.Source, error) {
		return &stubSource{baseURL: cfg.HttpFetcher.URL}, nil
	})

	builder := NewBuilder(newTestDB(t))
	recipeRuntime, err := builder.BuildRecipe(context.Background(), &dao.CustomRecipeV2{
		ID:         "recipe-source-input-spec",
		Craft:      "proxy",
		SourceType: string(testSourceType),
		SourceConfig: `{
			"http_fetcher":{"url":"https://example.com/feed.xml"}
		}`,
	})
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/feed.xml", recipeRuntime.BaseURL)
	assert.Equal(t, string(testSourceType), recipeRuntime.SourceType)

	feed, err := recipeRuntime.Fetch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "stub-feed", feed.Title)
}

func TestProxyRecipeFetch_UsesDefaultUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Proxy Feed</title>
    <link>https://example.com/</link>
    <description>Proxy test feed</description>
    <item>
      <title>Item 1</title>
      <link>https://example.com/item-1</link>
      <description>Hello</description>
    </item>
  </channel>
</rss>`)
	}))
	defer server.Close()

	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.CustomRecipeV2{
		ID:         "proxy-runtime-default-ua",
		Craft:      "proxy",
		SourceType: string(constant.SourceRSS),
		SourceConfig: `{
			"type":"rss",
			"http_fetcher":{"url":"` + server.URL + `"}
		}`,
	}).Error)

	builder := NewBuilder(db)
	provider, err := builder.BuildRecipeProvider(context.Background(), "proxy-runtime-default-ua")
	require.NoError(t, err)

	feed, err := provider.Fetch(context.Background())
	require.NoError(t, err)
	require.NotNil(t, feed)
	assert.Equal(t, "FeedCraft/2.0", gotUA)
}

func TestBuildAggregator(t *testing.T) {
	aggregator, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "sort", Option: map[string]string{"by": "date_desc"}},
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_link"}},
		{Type: "limit", Option: map[string]string{"max": "2"}},
	})
	require.NoError(t, err)
	require.NotNil(t, aggregator)

	now := time.Now()
	result, err := aggregator(context.Background(), &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Link: "http://a.com", Updated: now.Add(-1 * time.Hour)},
			{Id: "2", Link: "http://b.com", Updated: now},
			{Id: "3", Link: "http://a.com", Updated: now.Add(1 * time.Hour)},
			{Id: "4", Link: "http://c.com", Updated: now.Add(-2 * time.Hour)},
		},
	})
	require.NoError(t, err)
	require.Len(t, result.Articles, 2)
	assert.Equal(t, "3", result.Articles[0].Id)
	assert.Equal(t, "2", result.Articles[1].Id)
}

func TestBuildAggregator_InvalidLimit(t *testing.T) {
	_, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "limit", Option: map[string]string{"max": "0"}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid max")
}

func TestDeduplicate_ByTitle(t *testing.T) {
	aggregator, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_title"}},
	})
	require.NoError(t, err)
	require.NotNil(t, aggregator)

	result, err := aggregator(context.Background(), &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Title: "Hello World", Link: "http://a.com/1"},
			{Id: "2", Title: "Hello World", Link: "http://b.com/2"},  // duplicate title
			{Id: "3", Title: "Hello World!", Link: "http://c.com/3"}, // different (extra !)
			{Id: "4", Title: "", Link: "http://d.com/4"},             // empty title: kept unconditionally
			{Id: "5", Title: "", Link: "http://e.com/5"},             // another empty title: also kept
		},
	})
	require.NoError(t, err)
	// "Hello World" kept once, "Hello World!" kept, both empty-title articles kept
	require.Len(t, result.Articles, 4)
	assert.Equal(t, "1", result.Articles[0].Id)
	assert.Equal(t, "3", result.Articles[1].Id)
	assert.Equal(t, "4", result.Articles[2].Id)
	assert.Equal(t, "5", result.Articles[3].Id)
}

func TestDeduplicate_BySimhash_IdenticalContent(t *testing.T) {
	aggregator, err := BuildAggregator([]dao.AggregatorStep{
		// threshold=5 (normalized 0-100) → hamming ≈ 3
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_simhash", "threshold": "5"}},
	})
	require.NoError(t, err)
	require.NotNil(t, aggregator)

	text := "这是一篇关于人工智能的深度报道，探讨了大模型技术的最新进展。"
	result, err := aggregator(context.Background(), &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Title: text, Content: text},
			{Id: "2", Title: text, Content: text}, // exact duplicate
			{Id: "3", Title: "完全不同的文章：量子计算突破！", Content: "全新内容"},
		},
	})
	require.NoError(t, err)
	require.Len(t, result.Articles, 2)
	assert.Equal(t, "1", result.Articles[0].Id)
	assert.Equal(t, "3", result.Articles[1].Id)
}

func TestDeduplicate_BySimhash_InvalidThreshold(t *testing.T) {
	for _, bad := range []string{"999", "101", "-1", "abc"} {
		_, err := BuildAggregator([]dao.AggregatorStep{
			{Type: "deduplicate", Option: map[string]string{"strategy": "by_simhash", "threshold": bad}},
		})
		require.Error(t, err, "expected error for threshold=%q", bad)
		assert.Contains(t, err.Error(), "threshold", "expected 'threshold' in error for threshold=%q", bad)
	}
}

func TestDeduplicate_BySimhash_DefaultThreshold(t *testing.T) {
	// Omitting threshold should use default (5, normalized 0-100) without error
	aggregator, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_simhash"}},
	})
	require.NoError(t, err)
	require.NotNil(t, aggregator)

	result, err := aggregator(context.Background(), &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Title: "test"},
			{Id: "2", Title: "unique content here"},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, result.Articles)
}

func TestDeduplicate_BySimhash_BoundaryThresholds(t *testing.T) {
	// threshold=0 and threshold=100 are both valid
	for _, v := range []string{"0", "100"} {
		agg, err := BuildAggregator([]dao.AggregatorStep{
			{Type: "deduplicate", Option: map[string]string{"strategy": "by_simhash", "threshold": v}},
		})
		require.NoError(t, err, "threshold=%s should be valid", v)
		require.NotNil(t, agg)
	}
}

func TestDeduplicate_ByEmbedding_InvalidStrategy(t *testing.T) {
	_, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_unknown"}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid strategy")
}

func TestDeduplicate_ByEmbedding_ThresholdOutOfRange(t *testing.T) {
	_, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_embedding", "threshold": "1.5"}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "threshold")
}

func TestBuildTopicProvider_NestedTopics(t *testing.T) {
	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.TopicFeed{
		ID:        "child",
		Title:     "Child Topic",
		InputURIs: []string{"https://example.com/feed.xml"},
	}).Error)
	require.NoError(t, db.Create(&dao.TopicFeed{
		ID:        "parent",
		Title:     "Parent Topic",
		InputURIs: []string{"feedcraft://topic/child"},
		AggregatorConfig: []dao.AggregatorStep{
			{Type: "limit", Option: map[string]string{"max": "5"}},
		},
	}).Error)

	builder := NewBuilder(db)
	provider, err := builder.BuildTopicProvider(context.Background(), "parent")
	require.NoError(t, err)

	topicProvider, ok := provider.(*engine.TopicFeed)
	require.True(t, ok)
	assert.Equal(t, "parent", topicProvider.ID)
	assert.Len(t, topicProvider.Inputs, 1)
	_, ok = topicProvider.Inputs[0].(*engine.TopicFeed)
	assert.True(t, ok)
	require.NotNil(t, topicProvider.Aggregator)

	result, err := topicProvider.Aggregator(context.Background(), &model.CraftFeed{
		Articles: []*model.CraftArticle{{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"}, {Id: "5"}, {Id: "6"}},
	})
	require.NoError(t, err)
	require.Len(t, result.Articles, 5)
}

func TestBuildTopicProvider_CycleDetection(t *testing.T) {
	db := newTestDB(t)
	require.NoError(t, db.Create(&dao.TopicFeed{
		ID:        "A",
		InputURIs: []string{"feedcraft://topic/B"},
	}).Error)
	require.NoError(t, db.Create(&dao.TopicFeed{
		ID:        "B",
		InputURIs: []string{"feedcraft://topic/A"},
	}).Error)

	builder := NewBuilder(db)
	_, err := builder.BuildTopicProvider(context.Background(), "A")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "topic dependency cycle detected: A -> B -> A")
}

func TestBuildProviderFromInput_MissingRecipeFailsEarly(t *testing.T) {
	builder := NewBuilder(newTestDB(t))
	_, err := builder.BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindURI,
		URI:  "feedcraft://recipe/missing",
	}, nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}, &dao.TopicFeed{}, &dao.CraftAtom{}))
	return db
}

type stubSource struct {
	baseURL          string
	itemLinkOverride string
}

func (s *stubSource) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	now := time.Now()
	return &model.CraftFeed{
		Title:   "stub-feed",
		Link:    s.baseURL,
		Id:      s.baseURL + "/rss",
		Created: now,
		Updated: now,
		Articles: []*model.CraftArticle{
			{
				Title:   "Item 1",
				Link:    firstNonEmpty(s.itemLinkOverride, s.baseURL+"/item-1"),
				Id:      "item-1",
				Created: now,
				Updated: now,
			},
		},
	}, nil
}

func (s *stubSource) BaseURL() string {
	return s.baseURL
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func registerTestSource(t *testing.T, sourceType constant.SourceType, factory source.SourceFactory) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("register test source %q panicked: %v", sourceType, r)
		}
	}()
	source.Register(sourceType, factory)
}
