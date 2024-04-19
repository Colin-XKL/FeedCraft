package recipe

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
	"github.com/samber/lo"
	"time"
)

func convertEnclosure(item *gofeed.Enclosure, index int) feeds.Enclosure {
	return feeds.Enclosure{
		Url:    item.URL,
		Length: item.Length,
		Type:   item.Type,
	}
}

func ProxyFeed(c *gin.Context) {
	// Handle request for version 2 of users route
	fp := gofeed.NewParser()
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(500, "empty feed url")
		return
	}
	parsedFeed, _ := fp.ParseURL(feedUrl)
	fmt.Println(parsedFeed.Title)
	ret := feeds.Feed{
		Title: parsedFeed.Title,
		Link: &feeds.Link{
			Href: parsedFeed.Link,
		},
		Description: parsedFeed.Description,
		Author: &feeds.Author{
			Name:  parsedFeed.Author.Name,
			Email: parsedFeed.Author.Email,
		},
		Updated: *parsedFeed.UpdatedParsed,
		Created: *parsedFeed.PublishedParsed,
		Id:      parsedFeed.FeedLink,
		Items: lo.Map(parsedFeed.Items, func(item *gofeed.Item, index int) *feeds.Item {

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
			return &feeds.Item{
				Title: item.Title,
				Link: &feeds.Link{
					Href: item.Link,
				},
				Author: &feeds.Author{
					Name:  item.Author.Name,
					Email: item.Author.Email,
				},
				Description: item.Description,
				Id:          item.GUID,
				Updated:     updatedTime,
				Created:     publishedTime,
				Content:     item.Content,
			}

		}),
		Copyright: parsedFeed.Copyright,
		Image: &feeds.Image{
			Url:   parsedFeed.Image.URL,
			Title: parsedFeed.Image.Title,
			Link:  parsedFeed.Image.URL,
		},
	}

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
