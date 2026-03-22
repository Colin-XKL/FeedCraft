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

// FeedProcessor represents any node that takes a CraftFeed, processes it, and returns a new CraftFeed.
// Examples: AtomCrafts (Translate, FullText, Summary), FlowCraft, Aggregator.
type FeedProcessor interface {
	Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error)
}
