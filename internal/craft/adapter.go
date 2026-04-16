package craft

import (
	"context"

	"FeedCraft/internal/model"
)

// LegacyOptionAdapter wraps a legacy CraftOption so it can be used as a FeedProcessor.
type LegacyOptionAdapter struct {
	Option CraftOption
	Extra  ExtraPayload
}

// Process implements the engine.FeedProcessor interface
func (a *LegacyOptionAdapter) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	// 1. Convert internal model back to Gorilla feeds format for legacy processor
	legacyFeed := feed.ToFeedsFeed()

	// 2. Apply the legacy functional mutation
	if err := a.Option(legacyFeed, a.Extra); err != nil {
		return nil, err
	}

	// 3. Convert the mutated feeds.Feed back into the internal CraftFeed model
	return model.FromFeedsFeed(legacyFeed), nil
}
