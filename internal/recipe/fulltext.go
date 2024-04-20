package recipe

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const DefaultTimeout = 30 * time.Second

func GetCacheKeyForWebContent(url string) string {
	return fmt.Sprintf("%s_%s", constant.PrefixWebContent, url)
}

func extractFulltext(item *gofeed.Item) string {
	url := item.Link
	log.Printf("extract fulltext for url %s", url)

	articleContent := ""

	content, err := util.CacheGetString(GetCacheKeyForWebContent(url))
	if err != nil || content == "" {
		article, err := readability.FromURL(url, DefaultTimeout)
		if err != nil {
			logrus.Warning("failed to parse %s, %v\n", url, err)
		} else {
			articleContent = article.Content
			cacheErr := util.CacheSetString(GetCacheKeyForWebContent(url), article.Content, constant.WebContentExpire)
			if cacheErr != nil {
				logrus.Warnf("failed to cache %s, %v\n", url, cacheErr)
			}
		}
	} else {
		articleContent = content
	}
	return articleContent
}

func ExtractFulltextForFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(500, "empty feed url")
		return
	}
	fp := gofeed.NewParser()
	parsedFeed, _ := fp.ParseURL(feedUrl)

	ret := TransformFeed(parsedFeed, extractFulltext)

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
