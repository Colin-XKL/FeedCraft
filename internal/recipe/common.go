package recipe

import (
	"FeedCraft/internal/constant"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"time"
)

const DefaultTimeout = 30 * time.Second

func GetCacheKeyForWebContent(url string) string {
	return fmt.Sprintf("%s_%s", constant.PrefixWebContent, url)
}

type ContentTransformFunc func(item *gofeed.Item) string

func TransformFeed(parsedFeed *gofeed.Feed, transFunc ContentTransformFunc) feeds.Feed {
	updatedTimePointer := parsedFeed.UpdatedParsed
	updatedTime := time.Now()
	if updatedTimePointer != nil {
		updatedTime = *updatedTimePointer
	}

	publishedTimePointer := parsedFeed.PublishedParsed
	publishedTime := time.Now()
	if publishedTimePointer != nil {
		publishedTime = *publishedTimePointer
	}

	extractIterator := func(item *gofeed.Item, index int) *feeds.Item {
		return TransformFeedItem(item, transFunc)
	}

	ret := feeds.Feed{
		Title: parsedFeed.Title,
		Link: &feeds.Link{
			Href: parsedFeed.Link,
		},
		Description: parsedFeed.Description,
		Updated:     updatedTime,
		Created:     publishedTime,
		Id:          parsedFeed.FeedLink,
		Items:       lo.Map(parsedFeed.Items, extractIterator),
		Copyright:   parsedFeed.Copyright,
	}
	if parsedFeed.Author != nil {
		ret.Author = &feeds.Author{
			Name:  parsedFeed.Author.Name,
			Email: parsedFeed.Author.Email,
		}
	}
	if parsedFeed.Image != nil {
		ret.Image = &feeds.Image{
			Url:   parsedFeed.Image.URL,
			Title: parsedFeed.Image.Title,
			Link:  parsedFeed.Image.URL,
		}
	}
	return ret
}

func TransformFeedItem(item *gofeed.Item, transFunc func(item *gofeed.Item) string) *feeds.Item {
	updatedTimePointer := item.UpdatedParsed
	updatedTime := time.Now()
	if updatedTimePointer != nil {
		updatedTime = *updatedTimePointer
	}

	publishedTimePointer := item.PublishedParsed
	publishedTime := time.Now()
	if publishedTimePointer != nil {
		publishedTime = *publishedTimePointer
	}

	articleContent := transFunc(item)

	retItem := feeds.Item{
		Title: item.Title,
		Link: &feeds.Link{
			Href: item.Link,
		},
		Description: item.Description,
		Id:          item.GUID,
		Updated:     updatedTime,
		Created:     publishedTime,
		Content:     articleContent,
	}

	if item.Author != nil {
		retItem.Author = &feeds.Author{
			Name:  item.Author.Name,
			Email: item.Author.Email,
		}
	}
	return &retItem
}
