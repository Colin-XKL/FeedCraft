package craft

import (
	"context"
	"testing"

	"FeedCraft/internal/model"

	"github.com/gorilla/feeds"
	"github.com/stretchr/testify/require"
)

type nilResultProcessor struct{}

func (nilResultProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	_ = ctx
	_ = feed
	return nil, nil
}

func TestApplyLocalProcessorToLegacyFeed_NilProcessorResultDoesNotPanic(t *testing.T) {
	feed := &feeds.Feed{
		Title: "keep original",
		Items: []*feeds.Item{{Title: "item"}},
	}

	err := applyLocalProcessorToLegacyFeed(context.Background(), nilResultProcessor{}, feed)
	require.Error(t, err)
	require.Contains(t, err.Error(), "nil feed")
	require.Equal(t, "keep original", feed.Title)
	require.Len(t, feed.Items, 1)
}
