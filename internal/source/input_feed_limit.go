package source

import (
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultInputFeedItemLimit = 30

func inputFeedItemLimit() int {
	raw := strings.TrimSpace(util.GetEnvClient().GetString("INPUT_FEED_ITEM_LIMIT"))
	if raw == "" {
		return defaultInputFeedItemLimit
	}

	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 0 {
		logrus.Warnf("invalid FC_INPUT_FEED_ITEM_LIMIT %q; using default %d", raw, defaultInputFeedItemLimit)
		return defaultInputFeedItemLimit
	}
	return limit
}

func applyInputFeedItemLimit(feed *model.CraftFeed) {
	limit := inputFeedItemLimit()
	if feed == nil || limit == 0 || len(feed.Articles) <= limit {
		return
	}

	sort.SliceStable(feed.Articles, func(i, j int) bool {
		return inputArticleTime(feed.Articles[i]).After(inputArticleTime(feed.Articles[j]))
	})
	feed.Articles = feed.Articles[:limit]
}

// ApplyInputFeedItemLimit keeps the most recent configured number of input articles.
func ApplyInputFeedItemLimit(feed *model.CraftFeed) {
	applyInputFeedItemLimit(feed)
}

func inputArticleTime(article *model.CraftArticle) time.Time {
	if article == nil {
		return time.Time{}
	}
	if !article.Created.IsZero() {
		return article.Created
	}
	return article.Updated
}
