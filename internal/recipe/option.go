package recipe

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
		return strings.Trim(item.Content, " ")
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
