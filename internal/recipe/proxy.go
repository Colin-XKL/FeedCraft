package recipe

import (
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"strings"
)

func byPass(item *gofeed.Item) string {
	return strings.Trim(item.Content, " ")
}

func ProxyFeed(c *gin.Context) {
	feedUrl, ok := c.GetQuery("input_url")
	if !ok || len(feedUrl) == 0 {
		c.String(500, "empty feed url")
		return
	}

	fp := gofeed.NewParser()
	parsedFeed, _ := fp.ParseURL(feedUrl)

	ret := TransformFeed(parsedFeed, byPass)

	rssStr, err := ret.ToRss()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.Header("Content-Type", "application/xml")
	c.String(200, rssStr)
}
