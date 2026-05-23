package source

import (
	"context"
	"fmt"
	"strings"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
)

func init() {
	Register(constant.SourceInbox, inboxSourceFactory)
}

type InboxSource struct {
	InboxID string
	Config  *config.SourceConfig
}

func inboxSourceFactory(cfg *config.SourceConfig) (Source, error) {
	if cfg.InboxSource == nil || cfg.InboxSource.InboxID == "" {
		return nil, fmt.Errorf("inbox_source config with inbox_id is required for inbox source")
	}
	return &InboxSource{
		InboxID: cfg.InboxSource.InboxID,
		Config:  cfg,
	}, nil
}

func (s *InboxSource) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	db := util.GetDatabase()

	inbox, err := dao.GetInboxByID(db, s.InboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbox %s: %w", s.InboxID, err)
	}

	items, err := dao.ListInboxItems(db, s.InboxID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inbox items: %w", err)
	}

	feed := &model.CraftFeed{
		Title:       inbox.Title,
		Description: inbox.Description,
		Id:          fmt.Sprintf("inbox:%s", s.InboxID),
		Link:        s.BaseURL(),
		Articles:    make([]*model.CraftArticle, 0, len(items)),
	}

	for _, item := range items {
		article := &model.CraftArticle{
			Title:       item.Title,
			Link:        item.URL,
			Content:     item.Content,
			Description: item.Summary,
			Id:          item.ItemID,
			AuthorName:  item.Author,
			Created:     item.PublishedAt,
			Updated:     item.PublishedAt, // inbox item currently doesn't mutate business fields, fallback to published
		}
		feed.Articles = append(feed.Articles, article)
	}

	return feed, nil
}

func (s *InboxSource) BaseURL() string {
	env := util.GetEnvClient()
	siteBaseURL := strings.TrimRight(env.GetString("SITE_BASE_URL"), "/")
	return fmt.Sprintf("%s/inbox/%s", siteBaseURL, s.InboxID)
}
