package feedruntime

import (
	"context"
	"errors"
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

func TestBuildProviderFromInput_InvalidURI(t *testing.T) {
	builder := NewBuilder(newTestDB(t))
	_, err := builder.BuildProviderFromInput(context.Background(), InputSpec{
		Kind: InputKindURI,
		URI:  "feedcraft://recipe",
	}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing resource id")
}

func TestBuildAggregator(t *testing.T) {
	processor, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "sort", Option: map[string]string{"by": "date_desc"}},
		{Type: "deduplicate", Option: map[string]string{"strategy": "by_link"}},
		{Type: "limit", Option: map[string]string{"max": "10"}},
	})
	require.NoError(t, err)

	flow, ok := processor.(*engine.FlowCraftProcessor)
	require.True(t, ok)
	assert.Len(t, flow.Processors, 3)
}

func TestBuildAggregator_InvalidLimit(t *testing.T) {
	_, err := BuildAggregator([]dao.AggregatorStep{
		{Type: "limit", Option: map[string]string{"max": "0"}},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid max")
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
	assert.IsType(t, &engine.FlowCraftProcessor{}, topicProvider.Aggregator)
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
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}, &dao.TopicFeed{}))
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
