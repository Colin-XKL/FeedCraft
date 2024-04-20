package recipe

import (
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
			retItem := feeds.Item{
				Title: item.Title,
				Link: &feeds.Link{
					Href: item.Link,
				},
				Description: item.Description,
				Id:          item.GUID,
				Updated:     updatedTime,
				Created:     publishedTime,
				Content:     item.Content,
			}
			if item.Author != nil {
				retItem.Author = &feeds.Author{
					Name:  item.Author.Name,
					Email: item.Author.Email,
				}
			}
			return &retItem

		}),
		Copyright: parsedFeed.Copyright,
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
