package source

import (
	"FeedCraft/internal/config"
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

	assert.Equal(t, "https://www.google.com/s2/favicons?domain_url=https%3A%2F%2Fexample.com&sz=64", feed.ImageURL)
	assert.Equal(t, "Example Site", feed.ImageTitle)
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
