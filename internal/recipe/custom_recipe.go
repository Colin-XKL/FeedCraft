package recipe

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/craft"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/source"
	"FeedCraft/internal/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Scheduler *util.PreheatingScheduler

// ProcessRecipeByID encapsulates the core logic for fetching, generating, and crafting a recipe.
// It is designed to be reusable by both API handlers and background tasks like preheating.
func ProcessRecipeByID(ctx context.Context, recipeId string) (*feeds.Feed, error) {
	db := util.GetDatabase()
	channel, err := dao.GetChannelByID(db, recipeId)
	if err != nil {
		return nil, err
	}

	// 1. Parse SourceConfig to get the source of the feed
	var sourceConfig config.SourceConfig
	if err := json.Unmarshal([]byte(channel.SourceConfig), &sourceConfig); err != nil {
		return nil, fmt.Errorf("invalid source config: %w", err)
	}

	// ... (rest of the function uses sourceConfig)

	// 2. Get factory from registry
	factory, err := source.Get(sourceConfig.Type)
	if err != nil {
		return nil, err
	}

	// 3. Create source instance
	sourceInstance, err := factory(&sourceConfig)
	if err != nil {
		return nil, err
	}

	// 4. Generate the base feed
	baseFeed, err := sourceInstance.Generate(ctx)
	if err != nil {
		return nil, errors.New("failed to generate base feed: " + err.Error())
	}

	// 5. Get the base URL from the source for relative link resolution
	feedURL := sourceInstance.BaseURL()

	// 6. Process the feed through the craft flow
	processedFeed, err := craft.ProcessFeed(baseFeed, feedURL, channel.ProcessorName)
	if err != nil {
		return nil, errors.New("failed to process feed: " + err.Error())
	}

	return processedFeed, nil
}

// CustomRecipe is the Gin handler for serving a crafted recipe feed.
func CustomRecipe(c *gin.Context) {
	recipeId := c.Param("id")
	logrus.Infof("获取到的 recipe ID: [%s]", recipeId)

	processedFeed, err := ProcessRecipeByID(c.Request.Context(), recipeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// Schedule preheating
	logrus.Infof("add preheating task for recipe [%s]", recipeId)
	Scheduler.ScheduleTask(recipeId)

	// Render the final RSS
	rssStr, err := processedFeed.ToRss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to render RSS: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rssStr))
}