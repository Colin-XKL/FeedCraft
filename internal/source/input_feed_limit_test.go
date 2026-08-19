package source

import (
	"FeedCraft/internal/model"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInputFeedItemLimit(t *testing.T) {
	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "")
	assert.Equal(t, 30, inputFeedItemLimit())

	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "0")
	assert.Equal(t, 0, inputFeedItemLimit())

	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "12")
	assert.Equal(t, 12, inputFeedItemLimit())

	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "invalid")
	assert.Equal(t, 30, inputFeedItemLimit())

	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "-1")
	assert.Equal(t, 30, inputFeedItemLimit())
}

func TestApplyInputFeedItemLimit_SortsAndTruncatesBeforeProcessing(t *testing.T) {
	now := time.Now()
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "oldest", Created: now.Add(-3 * time.Hour)},
			{Id: "newest", Created: now},
			{Id: "updated-fallback", Updated: now.Add(-1 * time.Hour)},
			{Id: "undated"},
		},
	}
	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "2")

	applyInputFeedItemLimit(feed)

	assert.Equal(t, []string{"newest", "updated-fallback"}, articleIDs(feed.Articles))
}

func TestApplyInputFeedItemLimit_ZeroKeepsAllArticles(t *testing.T) {
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "first"},
			{Id: "second"},
		},
	}
	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "0")

	applyInputFeedItemLimit(feed)

	assert.Equal(t, []string{"first", "second"}, articleIDs(feed.Articles))
}

func TestPipelineSourceFetch_AppliesInputFeedItemLimit(t *testing.T) {
	now := time.Now()
	pipeline := &PipelineSource{
		Fetcher: inputLimitTestFetcher{},
		Parser: inputLimitTestParser{feed: &model.CraftFeed{
			Articles: []*model.CraftArticle{
				{Id: "older", Created: now.Add(-time.Hour)},
				{Id: "newer", Created: now},
			},
		}},
	}
	t.Setenv("FC_INPUT_FEED_ITEM_LIMIT", "1")

	feed, err := pipeline.Fetch(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []string{"newer"}, articleIDs(feed.Articles))
}

func articleIDs(articles []*model.CraftArticle) []string {
	ids := make([]string, 0, len(articles))
	for _, article := range articles {
		ids = append(ids, article.Id)
	}
	return ids
}

type inputLimitTestFetcher struct{}

func (inputLimitTestFetcher) Fetch(context.Context) ([]byte, error) {
	return []byte("feed"), nil
}

func (inputLimitTestFetcher) BaseURL() string {
	return ""
}

type inputLimitTestParser struct {
	feed *model.CraftFeed
}

func (p inputLimitTestParser) Parse([]byte) (*model.CraftFeed, error) {
	return p.feed, nil
}
