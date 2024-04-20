package recipe

import (
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const DefaultTimeout = 30 * time.Second

func GetCacheKeyForWebContent(url string) string {
	return fmt.Sprintf("%s_%s", constant.PrefixWebContent, url)
}

func extract(item *gofeed.Item, index int) *feeds.Item {

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

	//else {
	//	articleContent = article.Content
	//}

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

func ExtractFulltextForFeed(c *gin.Context) {
	fp := gofeed.NewParser()
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(500, "empty feed url")
		return
	}
	parsedFeed, _ := fp.ParseURL(feedUrl)

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

	ret := feeds.Feed{
		Title: parsedFeed.Title,
		Link: &feeds.Link{
			Href: parsedFeed.Link,
		},
		Description: parsedFeed.Description,
		Updated:     updatedTime,
		Created:     publishedTime,
		Id:          parsedFeed.FeedLink,
		Items:       lo.Map(parsedFeed.Items, extract),
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

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
