package craft

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
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

	fp := gofeed.NewParser()
	parsedFeed, err := fp.ParseURL(feedUrl)
	if err != nil {
		return ingredient, err
	}
	ingredient.parsedFeed = parsedFeed

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
			logrus.Errorf("failed to process article content for item [%s]. err: %v", item.Title, err)
			item.Content = fmt.Sprintf("Error processing content: %T", err)
			item.Description = item.Content
			return nil
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
			logrus.Errorf("failed to process article title for item [%s]. err: %v", item.Title, err)
			item.Title = fmt.Sprintf("Error processing title: %T", err)
			return nil
		}
		item.Title = transformed
		return nil
	}
}

type FeedItemProcessor func(feedItem *feeds.Item, payload ExtraPayload) error // 对每个feed item要执行的操作

// OptionTransformFeedItem 通用的feed item 处理
func OptionTransformFeedItem(processor FeedItemProcessor) CraftOption {
	return func(feed *feeds.Feed, payload ExtraPayload) error {
		for _, itemPointer := range feed.Items {
			err := processor(itemPointer, payload)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
