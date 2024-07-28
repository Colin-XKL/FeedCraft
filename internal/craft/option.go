package craft

import (
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"strings"
)

type FeedCraftIngredient struct {
	originalFeedUrl string
	parsedFeed      *gofeed.Feed
	OutputFeed      *feeds.Feed
}
type CraftOption func(*feeds.Feed) error

func NewCraftedFeedFromUrl(feedUrl string, options ...CraftOption) (FeedCraftIngredient, error) {
	ingredient := FeedCraftIngredient{originalFeedUrl: feedUrl}

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
	outputFeed := TransformFeed(parsedFeed, byPass)

	for _, option := range options {
		optionErr := option(&outputFeed)
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
	return func(item *feeds.Item) error {
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
	return func(item *feeds.Item) error {
		transformed, err := transFunc(item)
		if err != nil {
			return err
		}
		item.Title = transformed
		return nil
	}
}

type FeedItemProcessor func(feedItem *feeds.Item) error // 对每个feed item要执行的操作

// OptionTransformFeedItem 通用的feed item 处理
func OptionTransformFeedItem(processor FeedItemProcessor) CraftOption {
	return func(feed *feeds.Feed) error {
		for _, itemPointer := range feed.Items {
			err := processor(itemPointer)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
