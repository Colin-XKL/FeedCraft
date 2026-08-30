package craft

import (
	"context"

	"FeedCraft/internal/model"

	"github.com/gorilla/feeds"
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

func applyLocalProcessorToLegacyFeed(ctx context.Context, processor localProcessor, feed *feeds.Feed) error {
	if feed == nil {
		return nil
	}
	out, err := processor.Process(ctx, model.FromFeedsFeed(feed))
	if err != nil {
		return err
	}
	converted := out.ToFeedsFeed()
	*feed = *converted
	return nil
}

func articleFromFeedItem(item *feeds.Item) *model.CraftArticle {
	if item == nil {
		return nil
	}
	article := &model.CraftArticle{
		Title:       item.Title,
		Description: item.Description,
		Id:          item.Id,
		Updated:     item.Updated,
		Created:     item.Created,
		Content:     item.Content,
	}
	if item.Link != nil {
		article.Link = item.Link.Href
	}
	if item.Author != nil {
		article.AuthorName = item.Author.Name
		article.AuthorEmail = item.Author.Email
	}
	return article
}
