package engine

import (
	"context"

	"FeedCraft/internal/model"
)

// FeedProvider represents any node that can generate or output a CraftFeed.
// Examples: RawFeed sources (HTML, RSS, Search), RecipeFeed, TopicFeed.
type FeedProvider interface {
	Fetch(ctx context.Context) (*model.CraftFeed, error)
}

// CraftOption represents a feed transformation closure in the runtime graph.
type CraftOption func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error)

func ComposeOptions(options ...CraftOption) CraftOption {
	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		current := feed
		var err error
		for _, option := range options {
			if option == nil {
				continue
			}
			current, err = option(ctx, current)
			if err != nil {
				return nil, err
			}
		}
		return current, nil
	}
}
