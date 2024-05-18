package recipe

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const DefaultTimeout = 30 * time.Second

//	func GetCacheKeyForWebContent(url string) string {
//		return fmt.Sprintf("%s_%s", constant.PrefixWebContent, url)
//	}

func getCacheKey(namespace, id string) string {
	return fmt.Sprintf("%s_%s_%s", constant.PrefixWebContent, namespace, id)
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
		return TransformArticleContent(item, transFunc)
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

func CommonCraftHandlerUsingCraftOptionList(c *gin.Context, optionList []CraftOption) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(400, "empty feed url")
		return
	}

	craftedFeed, err := NewCraftedFeedFromUrl(feedUrl,
		optionList...,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	rssStr, err := craftedFeed.OutputFeed.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}

type RawTransformer func(item *feeds.Item) (string, error)

func GetCommonCachedTransformer(cacheKeyGenerator ContentCacheKeyGenerator, rawTransformer TransFunc, craftName string) TransFunc {
	ret := func(item *feeds.Item) (string, error) {
		originalTitle := item.Title
		logrus.Info("applying craft [%s] to article %s", craftName, originalTitle)

		final := ""
		hashVal, _ := cacheKeyGenerator(item)

		cached, err := util.CacheGetString(getCacheKey(craftName, hashVal))
		if err != nil || cached == "" {
			translated, err := rawTransformer(item)
			if err != nil {
				logrus.Warnf("failed to apply craft [%s] for article [%s], %v\n", craftName, originalTitle, err)
				return "", err
			} else {
				final = translated
				cacheErr := util.CacheSetString(getCacheKey(craftName, hashVal), translated, constant.WebContentExpire)
				if cacheErr != nil {
					logrus.Warnf("failed to cache result of craft [%s] for article [%s], %v\n", craftName,
						originalTitle, cacheErr)
				}
			}
		} else {
			final = cached
		}
		return final, nil
	}
	return ret
}

func TransformArticleContent(item *gofeed.Item, transFunc func(item *gofeed.Item) string) *feeds.Item {
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
		Description: articleContent,
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
