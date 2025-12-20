package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
)

type SearchFetchReq struct {
	Query string `json:"query"`
}

func PreviewSearchRSS(c *gin.Context) {
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
			Query: req.Query,
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

	feed, err := src.Generate(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Generation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[[]*gofeed.Item]{
		StatusCode: 0,
		Data:       feed.Items,
	})
}
