package craft

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

const DefaultExtractFulltextTimeout = 30 * time.Second

//	func GetCacheKeyForWebContent(url string) string {
//		return fmt.Sprintf("%s_%s", constant.PrefixWebContent, url)
//	}

func getCraftCacheKey(namespace, id string) string {
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
		logrus.Infof("applying craft [%s] to article [%s]", craftName, originalTitle)

		hashVal, _ := cacheKeyGenerator(item)
		cacheKey := getCraftCacheKey(craftName, hashVal)

		valFunc := func() (string, error) {
			ret, err := rawTransformer(item)
			if err != nil {
				logrus.Warnf("failed to apply craft [%s] for article [%s], %v\n", craftName, originalTitle, err)
			}
			return ret, err
		}

		return util.CachedFunc(cacheKey, valFunc)
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

	authorItem := item.Author
	if item.Author == nil {
		authorItem = lo.FirstOrEmpty(item.Authors)
	}
	if authorItem != nil {
		retItem.Author = &feeds.Author{
			Name:  authorItem.Name,
			Email: authorItem.Email,
		}
	}
	return &retItem
}
