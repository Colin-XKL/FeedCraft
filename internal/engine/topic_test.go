package engine

import (
	"context"
	"errors"
	"testing"
	"time"

	"FeedCraft/internal/model"
	"github.com/stretchr/testify/assert"
)

// MockProvider is a simple FeedProvider for testing.
type MockProvider struct {
	Feed  *model.CraftFeed
	Err   error
	Delay time.Duration
}

func (m *MockProvider) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	return m.Feed, m.Err
}

func TestTopicFeed_Fetch_Success(t *testing.T) {
	updated1 := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	created1 := time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC)
	updated2 := time.Date(2026, 4, 2, 11, 0, 0, 0, time.UTC)
	created2 := time.Date(2026, 4, 2, 8, 0, 0, 0, time.UTC)
	updated3 := time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC)
	created3 := time.Date(2026, 4, 3, 7, 0, 0, 0, time.UTC)

	provider1 := &MockProvider{
		Feed: &model.CraftFeed{
			Articles: []*model.CraftArticle{
				{Id: "1", Title: "Article 1", Updated: updated1, Created: created1},
				{Id: "2", Title: "Article 2", Updated: updated2, Created: created2},
			},
		},
	}
	provider2 := &MockProvider{
		Feed: &model.CraftFeed{
			Articles: []*model.CraftArticle{
				{Id: "3", Title: "Article 3", Updated: updated3, Created: created3},
			},
		},
	}

	topic := &TopicFeed{
		ID:     "topic-1",
		Inputs: []FeedProvider{provider1, provider2},
	}

	result, err := topic.Fetch(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "topic-1", result.Id)

	// Should contain 3 articles total
	assert.Len(t, result.Articles, 3)
	assert.True(t, result.Updated.Equal(updated3))
	assert.True(t, result.Created.Equal(created3))
}

func TestTopicFeed_Fetch_PartialFailure(t *testing.T) {
	provider1 := &MockProvider{
		Feed: &model.CraftFeed{
			Articles: []*model.CraftArticle{
				{Id: "1", Title: "Article 1"},
			},
		},
	}
	providerFail := &MockProvider{
		Err: errors.New("network error"),
	}

	topic := &TopicFeed{
		ID:     "topic-2",
		Inputs: []FeedProvider{provider1, providerFail},
	}

	// Should not return error overall, just log the warning and return partial success
	result, err := topic.Fetch(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Should contain 1 article from provider1
	assert.Len(t, result.Articles, 1)
	assert.Equal(t, "1", result.Articles[0].Id)
}

func TestTopicFeed_Fetch_WithAggregator(t *testing.T) {
	updated1 := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	updated2 := time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC)
	updated3 := time.Date(2026, 4, 3, 10, 0, 0, 0, time.UTC)
	created1 := time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC)
	created2 := time.Date(2026, 4, 2, 9, 0, 0, 0, time.UTC)
	created3 := time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC)

	provider := &MockProvider{
		Feed: &model.CraftFeed{
			Articles: []*model.CraftArticle{
				{Id: "1", Title: "A", Updated: updated1, Created: created1},
				{Id: "2", Title: "B", Updated: updated2, Created: created2},
				{Id: "3", Title: "C", Updated: updated3, Created: created3},
			},
		},
	}

	// An aggregator that limits to 2 items
	aggregator := &LimitProcessor{MaxItems: 2}

	topic := &TopicFeed{
		Inputs:     []FeedProvider{provider},
		Aggregator: aggregator,
	}

	result, err := topic.Fetch(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Articles, 2)
	assert.True(t, result.Updated.Equal(updated2))
	assert.True(t, result.Created.Equal(created2))
}

func TestTopicFeed_Fetch_AllInputsFailed(t *testing.T) {
	providerFail1 := &MockProvider{Err: errors.New("network error")}
	providerFail2 := &MockProvider{Err: errors.New("timeout")}

	topic := &TopicFeed{
		ID:     "topic-3",
		Inputs: []FeedProvider{providerFail1, providerFail2},
	}

	result, err := topic.Fetch(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "all upstream providers failed")
}
