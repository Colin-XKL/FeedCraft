package source

import (
	"context"
	"fmt"
	"net/url"

	"FeedCraft/internal/model"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

// LegacySource keeps the old Source contract for compatibility during migration.
type LegacySource interface {
	Generate(ctx context.Context) (*gofeed.Feed, error)
	BaseURL() string
}

// LegacySourceAdapter wraps a legacy Source so it can be used as a FeedProvider.
type LegacySourceAdapter struct {
	LegacySource LegacySource
}

func (a *LegacySourceAdapter) Fetch(ctx context.Context) (*model.CraftFeed, error) {
	rawFeed, err := a.LegacySource.Generate(ctx)
	if err != nil {
		return nil, err
	}

	cf := model.FromGofeed(rawFeed)

	// Ensure feed link is absolute, replicating logic from legacy TransformFeed
	cf.Link = getAbsFeedLink(a.LegacySource.BaseURL(), cf.Link)

	return cf, nil
}

func getAbsFeedLink(feedUrl, feedLinkAttr string) string {
	feedLinkUrl, err := url.Parse(feedLinkAttr)
	if err != nil || feedLinkUrl == nil {
		logrus.Warnf("invalid feed link url [%s] for feed [%s]", feedLinkAttr, feedUrl)
	} else {
		if feedLinkUrl.IsAbs() {
			return feedLinkAttr
		}
	}
	parsedFeedUrl, err := url.Parse(feedUrl)
	if err != nil {
		logrus.Errorf("invalid feed url [%s]. err: %v", feedUrl, err)
	} else {
		if feedLinkAttr != "" && feedLinkUrl != nil {
			return parsedFeedUrl.ResolveReference(feedLinkUrl).String()
		}
		return fmt.Sprintf("%s://%s", parsedFeedUrl.Scheme, parsedFeedUrl.Host)
	}
	return feedLinkAttr
}
