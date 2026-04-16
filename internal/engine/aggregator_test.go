package engine

import (
	"context"
	"testing"
	"time"

	"FeedCraft/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicateProcessor(t *testing.T) {
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Link: "http://a.com"},
			{Id: "2", Link: "http://b.com"},
			{Id: "3", Link: "http://a.com"}, // Duplicate link
			{Id: "2", Link: "http://c.com"}, // Duplicate ID
		},
	}

	t.Run("By Link", func(t *testing.T) {
		p := &DeduplicateProcessor{Strategy: "by_link"}
		// Clone feed to avoid mutating original for next test
		f := &model.CraftFeed{Articles: append([]*model.CraftArticle{}, feed.Articles...)}

		res, err := p.Process(context.Background(), f)
		assert.NoError(t, err)
		assert.Len(t, res.Articles, 3)
		assert.Equal(t, "1", res.Articles[0].Id)
		assert.Equal(t, "2", res.Articles[1].Id)
		assert.Equal(t, "2", res.Articles[2].Id) // ID 2 is kept twice because links differ (b.com, c.com)
	})

	t.Run("By ID", func(t *testing.T) {
		p := &DeduplicateProcessor{Strategy: "by_id"}
		f := &model.CraftFeed{Articles: append([]*model.CraftArticle{}, feed.Articles...)}

		res, err := p.Process(context.Background(), f)
		assert.NoError(t, err)
		assert.Len(t, res.Articles, 3)
		assert.Equal(t, "http://a.com", res.Articles[0].Link)
		assert.Equal(t, "http://b.com", res.Articles[1].Link)
		assert.Equal(t, "http://a.com", res.Articles[2].Link) // Link a.com is kept twice because IDs differ (1, 3)
	})
}

func TestSortProcessor(t *testing.T) {
	now := time.Now()
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Updated: now.Add(-1 * time.Hour), QualityScore: 10}, // oldest
			{Id: "2", Updated: now, QualityScore: 30},                     // newest
			{Id: "3", Updated: now.Add(-30 * time.Minute), QualityScore: 20},
		},
	}

	t.Run("Date Descending", func(t *testing.T) {
		p := &SortProcessor{SortBy: "date_desc"}
		f := &model.CraftFeed{Articles: append([]*model.CraftArticle{}, feed.Articles...)}
		res, err := p.Process(context.Background(), f)
		assert.NoError(t, err)
		assert.Equal(t, "2", res.Articles[0].Id) // newest first
		assert.Equal(t, "3", res.Articles[1].Id)
		assert.Equal(t, "1", res.Articles[2].Id)
	})

	t.Run("Quality Descending", func(t *testing.T) {
		p := &SortProcessor{SortBy: "quality_desc"}
		f := &model.CraftFeed{Articles: append([]*model.CraftArticle{}, feed.Articles...)}
		res, err := p.Process(context.Background(), f)
		assert.NoError(t, err)
		assert.Equal(t, "2", res.Articles[0].Id) // highest quality first
		assert.Equal(t, "3", res.Articles[1].Id)
		assert.Equal(t, "1", res.Articles[2].Id)
	})
}

func TestLimitProcessor(t *testing.T) {
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1"}, {Id: "2"}, {Id: "3"}, {Id: "4"},
		},
	}

	p := &LimitProcessor{MaxItems: 2}
	res, err := p.Process(context.Background(), feed)
	assert.NoError(t, err)
	assert.Len(t, res.Articles, 2)
	assert.Equal(t, "1", res.Articles[0].Id)
	assert.Equal(t, "2", res.Articles[1].Id)
}

func TestFlowCraftProcessor_Aggregator(t *testing.T) {
	now := time.Now()
	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Link: "http://a.com", Updated: now.Add(-1 * time.Hour)},
			{Id: "2", Link: "http://b.com", Updated: now},
			{Id: "3", Link: "http://a.com", Updated: now.Add(1 * time.Hour)}, // duplicate link, but newest!
			{Id: "4", Link: "http://c.com", Updated: now.Add(-2 * time.Hour)},
		},
	}

	// Important: To keep the *newest* duplicate, we should sort BEFORE we deduplicate.
	aggregator := &FlowCraftProcessor{
		Processors: []FeedProcessor{
			&SortProcessor{SortBy: "date_desc"},
			&DeduplicateProcessor{Strategy: "by_link"},
			&LimitProcessor{MaxItems: 2},
		},
	}

	res, err := aggregator.Process(context.Background(), feed)
	assert.NoError(t, err)
	assert.Len(t, res.Articles, 2)

	// Because we sort desc first, the order is: 3, 2, 1, 4
	// Then we deduplicate: keep 3, keep 2, drop 1 (dup link of 3), keep 4
	// Then limit 2: keep 3, 2.
	assert.Equal(t, "3", res.Articles[0].Id)
	assert.Equal(t, "2", res.Articles[1].Id)
}
