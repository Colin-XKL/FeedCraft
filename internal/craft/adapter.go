package craft

import (
	"context"
	"fmt"

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
	if processor == nil {
		return fmt.Errorf("processor is nil")
	}
	originalItems := feed.Items
	out, err := processor.Process(ctx, model.FromFeedsFeed(feed))
	if err != nil {
		return err
	}
	converted := out.ToFeedsFeed()
	if converted == nil {
		return fmt.Errorf("processor returned nil feed")
	}
	restoreLegacyItemMetadata(originalItems, converted)
	*feed = *converted
	return nil
}

func restoreLegacyItemMetadata(originals []*feeds.Item, converted *feeds.Feed) {
	if converted == nil || len(originals) == 0 {
		return
	}
	index := make(map[string][]*feeds.Item, len(originals))
	for _, item := range originals {
		if item == nil {
			continue
		}
		key := legacyItemKey(item)
		index[key] = append(index[key], item)
	}
	for _, item := range converted.Items {
		if item == nil {
			continue
		}
		key := legacyItemKey(item)
		matches := index[key]
		if len(matches) == 0 {
			continue
		}
		orig := matches[0]
		index[key] = matches[1:]
		if item.Enclosure == nil {
			item.Enclosure = orig.Enclosure
		}
		if item.IsPermaLink == "" {
			item.IsPermaLink = orig.IsPermaLink
		}
	}
}

func legacyItemKey(item *feeds.Item) string {
	link := ""
	if item.Link != nil {
		link = item.Link.Href
	}
	if item.Id != "" {
		return item.Id + "\x00" + link
	}
	return item.Title + "\x00" + link
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
	if item.Source != nil {
		article.Source = item.Source.Href
	}
	return article
}
