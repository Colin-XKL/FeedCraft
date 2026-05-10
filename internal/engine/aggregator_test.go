package engine

import (
	"context"
	"testing"
	"time"

	"FeedCraft/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposeOptions_DeduplicateByLink(t *testing.T) {
	aggregator := ComposeOptions(
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			if feed == nil || len(feed.Articles) == 0 {
				return feed, nil
			}
			cloned := cloneFeed(feed)
			seen := make(map[string]bool)
			unique := make([]*model.CraftArticle, 0, len(cloned.Articles))
			for _, article := range cloned.Articles {
				if article == nil {
					continue
				}
				if article.Link == "" || !seen[article.Link] {
					if article.Link != "" {
						seen[article.Link] = true
					}
					unique = append(unique, article)
				}
			}
			cloned.Articles = unique
			return cloned, nil
		},
	)

	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Link: "http://a.com"},
			{Id: "2", Link: "http://b.com"},
			{Id: "3", Link: "http://a.com"},
			{Id: "2", Link: "http://c.com"},
		},
	}

	result, err := aggregator(context.Background(), feed)
	require.NoError(t, err)
	assert.Len(t, result.Articles, 3)
	assert.Equal(t, "1", result.Articles[0].Id)
	assert.Equal(t, "2", result.Articles[1].Id)
	assert.Equal(t, "2", result.Articles[2].Id)
	assert.Len(t, feed.Articles, 4)
}

func TestComposeOptions_SortAndLimit(t *testing.T) {
	now := time.Now()
	aggregator := ComposeOptions(
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			if feed == nil || len(feed.Articles) <= 1 {
				return feed, nil
			}
			cloned := cloneFeed(feed)
			sortArticlesByUpdatedDesc(cloned.Articles)
			return cloned, nil
		},
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			if feed == nil || len(feed.Articles) <= 2 {
				return feed, nil
			}
			cloned := cloneFeed(feed)
			cloned.Articles = cloned.Articles[:2]
			return cloned, nil
		},
	)

	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Updated: now.Add(-1 * time.Hour)},
			{Id: "2", Updated: now},
			{Id: "3", Updated: now.Add(-30 * time.Minute)},
		},
	}

	result, err := aggregator(context.Background(), feed)
	require.NoError(t, err)
	assert.Equal(t, "2", result.Articles[0].Id)
	assert.Equal(t, "3", result.Articles[1].Id)
	assert.Len(t, result.Articles, 2)
}

func TestComposeOptions_ComposedPipeline(t *testing.T) {
	now := time.Now()
	aggregator := ComposeOptions(
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			cloned := cloneFeed(feed)
			sortArticlesByUpdatedDesc(cloned.Articles)
			return cloned, nil
		},
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			cloned := cloneFeed(feed)
			seen := make(map[string]bool)
			unique := make([]*model.CraftArticle, 0, len(cloned.Articles))
			for _, article := range cloned.Articles {
				if article == nil {
					continue
				}
				if article.Link == "" || !seen[article.Link] {
					if article.Link != "" {
						seen[article.Link] = true
					}
					unique = append(unique, article)
				}
			}
			cloned.Articles = unique
			return cloned, nil
		},
		func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
			_ = ctx
			cloned := cloneFeed(feed)
			if len(cloned.Articles) > 2 {
				cloned.Articles = cloned.Articles[:2]
			}
			return cloned, nil
		},
	)

	feed := &model.CraftFeed{
		Articles: []*model.CraftArticle{
			{Id: "1", Link: "http://a.com", Updated: now.Add(-1 * time.Hour)},
			{Id: "2", Link: "http://b.com", Updated: now},
			{Id: "3", Link: "http://a.com", Updated: now.Add(1 * time.Hour)},
			{Id: "4", Link: "http://c.com", Updated: now.Add(-2 * time.Hour)},
		},
	}

	result, err := aggregator(context.Background(), feed)
	require.NoError(t, err)
	assert.Len(t, result.Articles, 2)
	assert.Equal(t, "3", result.Articles[0].Id)
	assert.Equal(t, "2", result.Articles[1].Id)
}

func cloneFeed(feed *model.CraftFeed) *model.CraftFeed {
	if feed == nil {
		return nil
	}
	cloned := *feed
	cloned.Articles = make([]*model.CraftArticle, 0, len(feed.Articles))
	for _, article := range feed.Articles {
		if article == nil {
			cloned.Articles = append(cloned.Articles, nil)
			continue
		}
		articleCopy := *article
		cloned.Articles = append(cloned.Articles, &articleCopy)
	}
	return &cloned
}

func sortArticlesByUpdatedDesc(articles []*model.CraftArticle) {
	for i := 0; i < len(articles); i++ {
		for j := i + 1; j < len(articles); j++ {
			if articles[j].Updated.After(articles[i].Updated) {
				articles[i], articles[j] = articles[j], articles[i]
			}
		}
	}
}
