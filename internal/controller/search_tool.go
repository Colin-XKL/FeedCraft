package controller

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/source/fetcher"
	"FeedCraft/internal/util"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchFetchReq struct {
	Query string `json:"query"`
}

func FetchSearch(c *gin.Context) {
	var req SearchFetchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: err.Error()})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{StatusCode: -1, Msg: "Query is required"})
		return
	}

	// Create a temporary fetcher
	f := &fetcher.SearchFetcher{
		Config: &config.SearchFetcherConfig{
			Query: req.Query,
		},
	}

	// Fetch data
	raw, err := f.Fetch(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, util.APIResponse[any]{StatusCode: -1, Msg: "Search failed: " + err.Error()})
		return
	}

	// Format JSON for display
	var prettyJSON bytes.Buffer
	if json.Unmarshal(raw, &struct{}{}) == nil {
		if err := json.Indent(&prettyJSON, raw, "", "  "); err == nil {
			c.JSON(http.StatusOK, util.APIResponse[string]{
				StatusCode: 0,
				Data:       prettyJSON.String(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, util.APIResponse[string]{
		StatusCode: 0,
		Data:       string(raw),
	})
}
