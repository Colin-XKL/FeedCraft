package recipe

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type FulltextExtractor func(url string, timeout time.Duration) (string, error)

func TrivialExtractor(url string, timeout time.Duration) (string, error) {
	article, err := readability.FromURL(url, timeout)
	return article.Content, err
}

func GetFulltextExtractor(extractor FulltextExtractor) ContentTransformFunc {
	extractFunc := func(item *gofeed.Item) string {
		url := item.Link
		log.Printf("extract fulltext for url %s", url)

		articleContent := ""

		cachedContent, err := util.CacheGetString(GetCacheKeyForWebContent(url))
		if err != nil || cachedContent == "" {
			articleStr, err := extractor(url, DefaultTimeout)
			if err != nil {
				logrus.Warnf("failed to parse %s, %v\n", url, err)
			} else {
				articleContent = articleStr
				cacheErr := util.CacheSetString(GetCacheKeyForWebContent(url), articleStr, constant.WebContentExpire)
				if cacheErr != nil {
					logrus.Warnf("failed to cache %s, %v\n", url, cacheErr)
				}
			}
		} else {
			articleContent = cachedContent
		}
		return articleContent
	}
	return extractFunc
}

func ExtractFulltextForFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(400, "empty feed url")
		return
	}
	fp := gofeed.NewParser()
	parsedFeed, _ := fp.ParseURL(feedUrl)

	ret := TransformFeed(parsedFeed, GetFulltextExtractor(TrivialExtractor))

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
