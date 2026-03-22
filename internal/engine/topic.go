package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"
	"FeedCraft/internal/observability"
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
	startedAt := time.Now()
	var g errgroup.Group
	var mu sync.Mutex
	var allArticles []*model.CraftArticle
	var failedInputs []map[string]any

	for _, input := range t.Inputs {
		// Capture loop variable
		provider := input
		g.Go(func() error {
			// We pass the context to allow cancellation or timeouts,
			// but we don't want one failure to cancel the others in an aggregator.
			feed, err := provider.Fetch(ctx)
			if err != nil {
				logrus.Warnf("TopicFeed [%s]: failed to fetch from a provider: %v", t.ID, err)
				providerType := fmt.Sprintf("%T", provider)
				mu.Lock()
				failedInputs = append(failedInputs, map[string]any{
					"provider_type": providerType,
					"error":         err.Error(),
					"error_kind":    observability.ClassifyError(err),
				})
				mu.Unlock()
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
		processedFeed, err := t.Aggregator.Process(ctx, mergedFeed)
		if err != nil {
			observability.Report(observability.ExecutionEvent{
				ResourceType: dao.ResourceTypeTopic,
				ResourceID:   t.ID,
				ResourceName: t.Title,
				Trigger:      observability.TriggerTopicAggregation,
				Status:       dao.ExecutionStatusFailure,
				ErrorKind:    observability.ClassifyError(err),
				Message:      err.Error(),
				Details: map[string]any{
					"failed_inputs": failedInputs,
				},
				RequestID: observability.RequestIDFromContext(ctx),
				Duration:  time.Since(startedAt),
			})
			return nil, err
		}
		reportTopicResult(ctx, t, processedFeed, failedInputs, startedAt)
		return processedFeed, nil
	}

	reportTopicResult(ctx, t, mergedFeed, failedInputs, startedAt)
	return mergedFeed, nil
}

func reportTopicResult(ctx context.Context, topic *TopicFeed, feed *model.CraftFeed, failedInputs []map[string]any, startedAt time.Time) {
	status := dao.ExecutionStatusSuccess
	errorKind := ""
	message := "topic executed successfully"
	if len(failedInputs) > 0 && len(feed.Articles) > 0 {
		status = dao.ExecutionStatusPartialSuccess
		errorKind = observability.ErrorKindUpstreamPartialFailure
		message = "topic completed with partial upstream failures"
	}
	if len(feed.Articles) == 0 && len(failedInputs) > 0 {
		status = dao.ExecutionStatusFailure
		errorKind = observability.ErrorKindEmptyFeed
		message = "topic failed because all upstream providers failed or produced no items"
	}

	observability.Report(observability.ExecutionEvent{
		ResourceType: dao.ResourceTypeTopic,
		ResourceID:   topic.ID,
		ResourceName: topic.Title,
		Trigger:      observability.TriggerTopicAggregation,
		Status:       status,
		ErrorKind:    errorKind,
		Message:      message,
		Details: map[string]any{
			"failed_inputs": failedInputs,
			"item_count":    len(feed.Articles),
		},
		RequestID: observability.RequestIDFromContext(ctx),
		Duration:  time.Since(startedAt),
	})
}
