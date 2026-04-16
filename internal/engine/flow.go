package engine

import (
	"context"

	"FeedCraft/internal/model"
)

// FlowCraftProcessor allows composing multiple FeedProcessors together sequentially.
type FlowCraftProcessor struct {
	Processors []FeedProcessor
}

func (f *FlowCraftProcessor) Process(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
	currentFeed := feed
	var err error

	for _, p := range f.Processors {
		currentFeed, err = p.Process(ctx, currentFeed)
		if err != nil {
			return nil, err
		}
	}

	return currentFeed, nil
}
