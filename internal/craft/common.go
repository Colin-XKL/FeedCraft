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
	"net/url"
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

func TransformFeed(parsedFeed *gofeed.Feed, feedUrl string, transFunc ContentTransformFunc) feeds.Feed {
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
			//Href: parsedFeed.Link,
			Href: getAbsFeedLink(feedUrl, parsedFeed.Link),
		},
		Description: parsedFeed.Description,
		Updated:     updatedTime,
		Created:     publishedTime,
		Id:          parsedFeed.FeedLink,
		Items:       lo.Map(parsedFeed.Items, extractIterator),
		Copyright:   parsedFeed.Copyright,
	}

	if len(parsedFeed.Authors) > 0 {
		ret.Author = &feeds.Author{
			Name:  parsedFeed.Authors[0].Name,
			Email: parsedFeed.Authors[0].Email,
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
				logrus.Warnf("failed to apply craft [%s] for article [%s], err: %v\n", craftName, originalTitle, err)
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

	authorItem := lo.FirstOrEmpty(item.Authors)
	if authorItem != nil {
		retItem.Author = &feeds.Author{
			Name:  authorItem.Name,
			Email: authorItem.Email,
		}
	}
	return &retItem
}

// 作用于 feed 级别, 确保获取到绝对url
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
		return fmt.Sprintf("%s://%s", parsedFeedUrl.Scheme, parsedFeedUrl.Host)
	}
	return feedLinkAttr
}

// 作用于 article级别, 确保获取到绝对路径url
// feedUrl: 原始feed文件的来源的url
// feedLinkAttr: feed 内容中的link字段
// feedItemUrl: feed 内容中每个文章的link字段
func getAbsLinkForFeedItem(feedUrl, feedLinkAttr, feedItemUrl string) string {
	feedLinkUrl, err := url.Parse(feedLinkAttr)
	if err != nil || feedLinkUrl == nil {
		logrus.Warnf("invalid feed link url [%s] for feed [%s]", feedLinkAttr, feedUrl)
	} else {
		if feedLinkUrl.IsAbs() {
			absFeedItemUrl, err := util.BuildAbsoluteURL(feedLinkAttr, feedItemUrl)
			if err != nil {
				logrus.Errorf("build absoluteURL failed. error: %v", err)
				return feedItemUrl
			}
			return absFeedItemUrl
		}
	}

	// if `link` attr in feed content is not an abs path, use feed url instead
	absFeedItemUrl, err := util.BuildAbsoluteURL(feedUrl, feedItemUrl)
	if err != nil {
		logrus.Errorf("build absoluteURL failed. error: %v", err)
		return feedItemUrl
	}
	return absFeedItemUrl
}
