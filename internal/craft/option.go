package craft

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/source/fetcher"
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"github.com/sirupsen/logrus"
	"strings"
)

type CraftedFeed struct {
	originalFeedUrl string
	parsedFeed      *gofeed.Feed
	OutputFeed      *feeds.Feed
}

// ExtraPayload extra info for crating feed
type ExtraPayload struct {
	originalFeedUrl string
}
type CraftOption func(*feeds.Feed, ExtraPayload) error

func NewCraftedFeedFromUrl(feedUrl string, options ...CraftOption) (CraftedFeed, error) {
	ingredient := CraftedFeed{originalFeedUrl: feedUrl}

	raw, err := (&fetcher.HttpFetcher{Config: &config.HttpFetcherConfig{
		URL:     feedUrl,
		Purpose: config.HttpFetcherPurposeFeed,
	}}).Fetch(context.Background())
	if err != nil {
		return ingredient, err
	}

	fp := gofeed.NewParser()
	parsedFeed, err := fp.Parse(bytes.NewReader(raw))
	if err != nil {
		return ingredient, err
	}
	ingredient.parsedFeed = parsedFeed

	return NewCraftedFeedFromGofeed(parsedFeed, feedUrl, options...)
}

func NewCraftedFeedFromGofeed(parsedFeed *gofeed.Feed, feedUrl string, options ...CraftOption) (CraftedFeed, error) {
	ingredient := CraftedFeed{originalFeedUrl: feedUrl, parsedFeed: parsedFeed}

	byPass := func(item *gofeed.Item) string {
		content := item.Content
		if len(content) == 0 {
			content = item.Description
		}
		return strings.Trim(content, " ")
	}
	outputFeed := TransformFeed(parsedFeed, feedUrl, byPass)

	payload := ExtraPayload{originalFeedUrl: feedUrl}
	for _, option := range options {
		optionErr := option(&outputFeed, payload)
		if optionErr != nil {
			return ingredient, optionErr
		}
	}

	ingredient.OutputFeed = &outputFeed
	return ingredient, nil
}

// TransFunc common transform func, such as translate the article content or title
type TransFunc func(item *feeds.Item) (string, error)

func GetArticleContentProcessor(transFunc TransFunc) FeedItemProcessor {
	return func(item *feeds.Item, payload ExtraPayload) error {
		transformed, err := transFunc(item)
		if err != nil {
			return err
		}
		item.Content = transformed
		item.Description = transformed
		return nil
	}
}

func GetArticleTitleProcessor(transFunc TransFunc) FeedItemProcessor {
	return func(item *feeds.Item, payload ExtraPayload) error {
		transformed, err := transFunc(item)
		if err != nil {
			return err
		}
		item.Title = transformed
		return nil
	}
}

type FeedItemProcessor func(feedItem *feeds.Item, payload ExtraPayload) error // 对每个feed item要执行的操作

// OptionTransformFeedItem 通用的feed item 处理
func OptionTransformFeedItem(processor FeedItemProcessor) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		if len(feed.Items) == 0 {
			return nil
		}

		// 1. 并发执行并将结果映射为错误切片
		errs := parallel.Map(feed.Items, func(item *feeds.Item, _ int) error {
			err := processor(item, payload)
			if err != nil {
				logrus.Warnf("failed to process item [%s], err: %v", item.Title, err)
			}
			return err
		})

		// 2. 检查是否全部失败
		if lo.EveryBy(errs, func(err error) bool { return err != nil }) {
			return fmt.Errorf("all items failed to process. last error: %v", lo.LastOrEmpty(errs))
		}

		return nil
	}
}
