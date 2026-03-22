package engine

import (
	"context"
	"sync"
	"time"

	"FeedCraft/internal/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// TopicFeed aggregates multiple FeedProviders into a single Feed.
type TopicFeed struct {
	ID          string
	Title       string
	Description string
	Link        string

	// Inputs are the upstream providers (RecipeFeeds, other TopicFeeds, or RawFeeds).
	Inputs []FeedProvider

	// Aggregator is an optional processor to handle deduplication, sorting, limiting, etc.
	Aggregator FeedProcessor
}

// Fetch implements the FeedProvider interface.
func (t *TopicFeed) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	var g errgroup.Group
	var mu sync.Mutex
	var allArticles []*model.CraftArticle

	for _, input := range t.Inputs {
		// Capture loop variable
		provider := input
		g.Go(func() error {
			// We pass the context to allow cancellation or timeouts,
			// but we don't want one failure to cancel the others in an aggregator.
			feed, err := provider.Fetch(ctx)
			if err != nil {
				logrus.Warnf("TopicFeed [%s]: failed to fetch from a provider: %v", t.ID, err)
				// Return nil so errgroup doesn't cancel other goroutines
				return nil
			}

			if feed != nil && len(feed.Articles) > 0 {
				mu.Lock()
				allArticles = append(allArticles, feed.Articles...)
				mu.Unlock()
			}
			return nil
		})
	}

	// We don't expect errors from Wait() since we return nil above, but we still wait for completion
	_ = g.Wait()

	mergedFeed := &model.CraftFeed{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Link:        t.Link,
		Updated:     time.Now(),
		Created:     time.Now(),
		Articles:    allArticles,
	}

	// If there's an aggregator pipeline (e.g., deduplicate -> sort -> limit), run it.
	if t.Aggregator != nil {
		return t.Aggregator.Process(ctx, mergedFeed)
	}

	return mergedFeed, nil
}
