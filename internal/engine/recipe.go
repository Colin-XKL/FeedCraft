package engine

import (
	"context"

	"FeedCraft/internal/model"
)

// RecipeFeed is a single-input feed node that applies an optional processor chain.
type RecipeFeed struct {
	ID          string
	Description string
	SourceType  string
	BaseURL     string
	CraftName   string

	Input     FeedProvider
	Processor FeedProcessor
}

// Fetch implements the FeedProvider interface.
func (r *RecipeFeed) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	if r.Input == nil {
		return nil, nil
	}

	rawFeed, err := r.Input.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	if r.Processor == nil {
		return rawFeed, nil
	}

	return r.Processor.Process(ctx, rawFeed)
}
