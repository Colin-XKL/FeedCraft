package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/model"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchFetchReq struct {
	Query        string `json:"query"`
	EnhancedMode bool   `json:"enhanced_mode"`
}

type SearchPreviewItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

func SearchPreview(c *gin.Context) {
	var req SearchFetchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: "Query is required"})
		return
	}

	cfg := &config.SourceConfig{
		Type: constant.SourceSearch,
		SearchFetcher: &config.SearchFetcherConfig{
			Query:        req.Query,
			EnhancedMode: req.EnhancedMode,
		},
	}

	factory, err := source.Get(constant.SourceSearch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{StatusCode: -1, Msg: "Factory not found: " + err.Error()})
		return
	}

	src, err := factory(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{StatusCode: -1, Msg: "Failed to create source: " + err.Error()})
		return
	}

	feed, err := src.Fetch(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{StatusCode: -1, Msg: "Generation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[[]SearchPreviewItem]{
		StatusCode: 0,
		Data:       buildSearchPreviewItems(feed),
	})
}

func buildSearchPreviewItems(feed *model.CraftFeed) []SearchPreviewItem {
	if feed == nil {
		return nil
	}

	items := make([]SearchPreviewItem, 0, len(feed.Articles))
	for _, article := range feed.Articles {
		if article == nil {
			continue
		}
		date := ""
		if !article.Created.IsZero() {
			date = article.Created.Format("2006-01-02 15:04:05")
		} else if !article.Updated.IsZero() {
			date = article.Updated.Format("2006-01-02 15:04:05")
		}

		items = append(items, SearchPreviewItem{
			Title:       article.Title,
			Link:        article.Link,
			Date:        date,
			Description: article.Description,
		})
	}
	return items
}
