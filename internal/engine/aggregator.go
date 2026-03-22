package engine

import (
	"context"
	"sort"
	"strings"

	"FeedCraft/internal/model"
)

// DeduplicateProcessor removes duplicate articles from a feed.
type DeduplicateProcessor struct {
	// Strategy defines how to identify duplicates. E.g., "by_link", "by_id". Defaults to "by_link".
	Strategy string
}

func (p *DeduplicateProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || len(feed.Articles) == 0 {
		return feed, nil
	}

	seen := make(map[string]bool)
	var uniqueArticles []*model.CraftArticle

	for _, article := range feed.Articles {
		var key string
		switch strings.ToLower(p.Strategy) {
		case "by_id":
			key = article.Id
		default: // default to by_link
			key = article.Link
		}

		if key == "" {
			// If key is empty, we don't deduplicate it (or we could choose to drop it, but keeping is safer)
			uniqueArticles = append(uniqueArticles, article)
			continue
		}

		if !seen[key] {
			seen[key] = true
			uniqueArticles = append(uniqueArticles, article)
		}
	}

	feed.Articles = uniqueArticles
	return feed, nil
}

// SortProcessor reorders the articles in a feed.
type SortProcessor struct {
	// SortBy defines the sorting criterion and order. E.g., "date_desc", "date_asc", "quality_desc".
	// Defaults to "date_desc".
	SortBy string
}

func (p *SortProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || len(feed.Articles) <= 1 {
		return feed, nil
	}

	articles := feed.Articles
	sortBy := strings.ToLower(p.SortBy)
	if sortBy == "" {
		sortBy = "date_desc" // default
	}

	sort.SliceStable(articles, func(i, j int) bool {
		a, b := articles[i], articles[j]
		switch sortBy {
		case "date_asc":
			return a.Updated.Before(b.Updated)
		case "quality_desc":
			return a.QualityScore > b.QualityScore
		case "quality_asc":
			return a.QualityScore < b.QualityScore
		default: // date_desc
			return a.Updated.After(b.Updated)
		}
	})

	feed.Articles = articles
	return feed, nil
}

// LimitProcessor truncates the number of articles in a feed.
type LimitProcessor struct {
	MaxItems int
}

func (p *LimitProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	if feed == nil || len(feed.Articles) == 0 {
		return feed, nil
	}

	if p.MaxItems > 0 && len(feed.Articles) > p.MaxItems {
		feed.Articles = feed.Articles[:p.MaxItems]
	}

	return feed, nil
}
