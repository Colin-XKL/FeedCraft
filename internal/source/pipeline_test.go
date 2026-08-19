package source

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/favicon"
	"FeedCraft/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipelineSourceApplyFeedIconSourceAuto(t *testing.T) {
	feed := &model.CraftFeed{
		Title:    "Example Site",
		Link:     "https://example.com/blog",
		ImageURL: "/assets/icon.png",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{IconSource: config.FeedIconSourceAuto},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://example.com/assets/icon.png", feed.ImageURL)
	assert.Equal(t, "Example Site", feed.ImageTitle)
}

func TestPipelineSourceApplyFeedIconSourceFaviconService(t *testing.T) {
	resetFaviconRegistry(t)
	feed := &model.CraftFeed{
		Title:    "Example Site",
		Link:     "https://example.com/blog",
		ImageURL: "https://example.com/custom.png",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{IconSource: config.FeedIconSourceFaviconService},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://t0.gstatic.cn/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https%3A%2F%2Fexample.com&size=64", feed.ImageURL)
	assert.Equal(t, "Example Site", feed.ImageTitle)
}

func TestPipelineSourceApplyFeedIconSourceUsesGlobalProvider(t *testing.T) {
	resetFaviconRegistry(t)
	if err := favicon.Replace(config.FaviconSettings{DefaultProviderID: favicon.ProviderYandex}); err != nil {
		t.Fatalf("replace favicon settings: %v", err)
	}
	feed := &model.CraftFeed{
		Title: "Example Site",
		Link:  "https://example.com/blog",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{IconSource: config.FeedIconSourceFaviconService},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://favicon.yandex.net/favicon/example.com", feed.ImageURL)
}

func TestPipelineSourceApplyFeedIconSourceRecipeProviderOverridesGlobal(t *testing.T) {
	resetFaviconRegistry(t)
	if err := favicon.Replace(config.FaviconSettings{DefaultProviderID: favicon.ProviderYandex}); err != nil {
		t.Fatalf("replace favicon settings: %v", err)
	}
	feed := &model.CraftFeed{
		Title: "Example Site",
		Link:  "https://example.com/blog",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{
				IconSource:      config.FeedIconSourceFaviconService,
				FaviconProvider: favicon.ProviderDuckDuckGo,
			},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://icons.duckduckgo.com/ip3/example.com.ico", feed.ImageURL)
}

func TestBuildFaviconServiceURLPreservesFullPageURLForCustomProvider(t *testing.T) {
	resetFaviconRegistry(t)
	if err := favicon.Replace(config.FaviconSettings{
		DefaultProviderID: "full_url",
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "full_url",
				Name:        "Full URL",
				URLTemplate: "https://icons.example.test/favicon?url={url_query}",
				Enabled:     true,
			},
		},
	}); err != nil {
		t.Fatalf("replace favicon settings: %v", err)
	}

	got := buildFaviconServiceURL("", "https://example.com/blog?page=2", "https://fallback.example.com")

	assert.Equal(t, "https://icons.example.test/favicon?url=https%3A%2F%2Fexample.com%2Fblog%3Fpage%3D2", got)
}

func TestPipelineSourceApplyFeedIconSourceAutoFallback(t *testing.T) {
	feed := &model.CraftFeed{
		Title: "Example Site",
		Link:  "https://example.com/blog",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{IconSource: config.FeedIconSourceAuto},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://example.com/favicon.ico", feed.ImageURL)
	assert.Equal(t, "Example Site", feed.ImageTitle)
}

func TestPipelineSourceApplyFeedIconSourceRejectsUnsafeExtractedIcon(t *testing.T) {
	feed := &model.CraftFeed{
		Title:    "Example Site",
		Link:     "https://example.com/blog",
		ImageURL: "javascript:alert(1)",
	}
	source := &PipelineSource{
		Config: &config.SourceConfig{
			FeedMeta: &config.FeedMetaConfig{IconSource: config.FeedIconSourceAuto},
		},
	}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Equal(t, "https://example.com/favicon.ico", feed.ImageURL)
	assert.Equal(t, "Example Site", feed.ImageTitle)
}

func TestPipelineSourceApplyFeedIconSourceDoesNotFallbackWithoutExplicitSource(t *testing.T) {
	feed := &model.CraftFeed{
		Title: "Example Site",
		Link:  "https://example.com/blog",
	}
	source := &PipelineSource{Config: &config.SourceConfig{}}

	source.applyFeedIconSource(feed, "https://example.com/blog/page")

	assert.Empty(t, feed.ImageURL)
	assert.Empty(t, feed.ImageTitle)
}

func resetFaviconRegistry(t *testing.T) {
	t.Helper()
	previous := favicon.Settings()
	t.Cleanup(func() {
		if err := favicon.Replace(previous); err != nil {
			t.Fatalf("restore favicon settings: %v", err)
		}
	})
	if err := favicon.Replace(favicon.DefaultSettings()); err != nil {
		t.Fatalf("reset favicon settings: %v", err)
	}
}
