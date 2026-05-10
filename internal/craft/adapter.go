package craft

import (
	"context"

	"FeedCraft/internal/model"
)

func AdaptLegacyOption(option LegacyCraftOption, extra ExtraPayload) CraftOption {
	if option == nil {
		return nil
	}

	return func(ctx context.Context, feed *model.CraftFeed) (*model.CraftFeed, error) {
		_ = ctx
		legacyFeed := feed.ToFeedsFeed()
		if err := option(legacyFeed, extra); err != nil {
			return nil, err
		}
		return model.FromFeedsFeed(legacyFeed), nil
	}
}

func AdaptLegacyOptions(options []LegacyCraftOption, extra ExtraPayload) CraftOption {
	if len(options) == 0 {
		return nil
	}

	adapted := make([]CraftOption, 0, len(options))
	for _, option := range options {
		adapted = append(adapted, AdaptLegacyOption(option, extra))
	}
	return ComposeOptions(adapted...)
}
