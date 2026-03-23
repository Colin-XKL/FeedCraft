package recipe

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/feedruntime"
	"FeedCraft/internal/observability"
	"FeedCraft/internal/util"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Scheduler *util.PreheatingScheduler

// ProcessRecipeByID encapsulates the core logic for fetching, generating, and crafting a recipe.
// It is designed to be reusable by both API handlers and background tasks like preheating.
func ProcessRecipeByID(ctx context.Context, recipeId string) (*feeds.Feed, error) {
	return ProcessRecipeByIDWithTrigger(ctx, recipeId, observability.TriggerUserRequest)
}

func ProcessRecipeByIDWithTrigger(ctx context.Context, recipeId string, trigger string) (*feeds.Feed, error) {
	startedAt := time.Now()
	db := util.GetDatabase()
	recipe, err := dao.GetCustomRecipeByIDV2(db, recipeId)
	if err != nil {
		observability.Report(observability.ExecutionEvent{
			ResourceType: dao.ResourceTypeRecipe,
			ResourceID:   recipeId,
			ResourceName: recipeId,
			Trigger:      trigger,
			Status:       dao.ExecutionStatusFailure,
			ErrorKind:    observability.ClassifyError(err),
			Message:      err.Error(),
			RequestID:    observability.RequestIDFromContext(ctx),
			Duration:     time.Since(startedAt),
		})
		return nil, err
	}

	recipeRuntime, err := feedruntime.NewBuilder(db).BuildRecipe(ctx, recipe)
	if err != nil {
		reportRecipeFailure(ctx, recipe, trigger, startedAt, err, nil)
		return nil, err
	}

	processedCraftFeed, err := recipeRuntime.Fetch(ctx)
	if err != nil {
		reportRecipeFailure(ctx, recipe, trigger, startedAt, err, map[string]any{
			"source_type": recipeRuntime.SourceType,
			"base_url":    recipeRuntime.BaseURL,
			"craft":       recipeRuntime.CraftName,
		})
		return nil, err
	}
	processedFeed := processedCraftFeed.ToFeedsFeed()

	observability.Report(observability.ExecutionEvent{
		ResourceType: dao.ResourceTypeRecipe,
		ResourceID:   recipe.ID,
		ResourceName: recipe.ID,
		Trigger:      trigger,
		Status:       dao.ExecutionStatusSuccess,
		Message:      fmt.Sprintf("recipe executed successfully with %d items", len(processedFeed.Items)),
		Details: map[string]any{
			"source_type": recipeRuntime.SourceType,
			"base_url":    recipeRuntime.BaseURL,
			"item_count":  len(processedFeed.Items),
		},
		RequestID: observability.RequestIDFromContext(ctx),
		Duration:  time.Since(startedAt),
	})

	return processedFeed, nil
}

// CustomRecipe is the Gin handler for serving a crafted recipe feed.
func CustomRecipe(c *gin.Context) {
	recipeId := c.Param("id")
	logrus.Infof("获取到的 recipe ID: [%s]", recipeId)
	requestID := fmt.Sprintf("recipe-%d", time.Now().UnixNano())
	ctx := observability.WithRequestID(c.Request.Context(), requestID)

	processedFeed, err := ProcessRecipeByID(ctx, recipeId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Recipe not found"})
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

func reportRecipeFailure(ctx context.Context, recipe *dao.CustomRecipeV2, trigger string, startedAt time.Time, err error, details map[string]any) {
	resourceName := ""
	resourceID := ""
	if recipe != nil {
		resourceName = recipe.ID
		resourceID = recipe.ID
	}
	observability.Report(observability.ExecutionEvent{
		ResourceType: dao.ResourceTypeRecipe,
		ResourceID:   resourceID,
		ResourceName: resourceName,
		Trigger:      trigger,
		Status:       dao.ExecutionStatusFailure,
		ErrorKind:    observability.ClassifyError(err),
		Message:      err.Error(),
		Details:      details,
		RequestID:    observability.RequestIDFromContext(ctx),
		Duration:     time.Since(startedAt),
	})
}
