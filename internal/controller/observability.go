package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/observability"
	"FeedCraft/internal/util"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
	"gorm.io/gorm"
)

type ResourceHealthView struct {
	ResourceType        string     `json:"resource_type"`
	ResourceID          string     `json:"resource_id"`
	ResourceName        string     `json:"resource_name"`
	CurrentStatus       string     `json:"current_status"`
	ConsecutiveFailures int        `json:"consecutive_failures"`
	LastSuccessAt       *time.Time `json:"last_success_at,omitempty"`
	LastFailureAt       *time.Time `json:"last_failure_at,omitempty"`
	LastErrorKind       string     `json:"last_error_kind,omitempty"`
	LastErrorMessage    string     `json:"last_error_message,omitempty"`
	PausedAt            *time.Time `json:"paused_at,omitempty"`
	PausedReason        string     `json:"paused_reason,omitempty"`
}

func ListObservableResources(c *gin.Context) {
	db := util.GetDatabase()
	resourceType := c.Query("resource_type")

	healthList, err := dao.ListResourceHealth(db, applyResourceHealthFilters(db, resourceType, c.Query("status")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	healthMap := make(map[string]*dao.ResourceHealth, len(healthList))
	for _, item := range healthList {
		key := item.ResourceType + ":" + item.ResourceID
		healthMap[key] = item
	}

	var result []ResourceHealthView
	if resourceType == "" || resourceType == dao.ResourceTypeRecipe {
		recipes, err := dao.ListCustomRecipeV2(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
		for _, recipe := range recipes {
			result = append(result, mergeResourceHealth(
				dao.ResourceTypeRecipe,
				recipe.ID,
				recipe.ID,
				healthMap[dao.ResourceTypeRecipe+":"+recipe.ID],
			))
		}
	}
	if resourceType == "" || resourceType == dao.ResourceTypeTopic {
		topics, err := dao.ListTopicFeeds(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
		for _, topic := range topics {
			result = append(result, mergeResourceHealth(
				dao.ResourceTypeTopic,
				topic.ID,
				topic.Title,
				healthMap[dao.ResourceTypeTopic+":"+topic.ID],
			))
		}
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: result})
}

func GetObservableResource(c *gin.Context) {
	db := util.GetDatabase()
	resourceType := c.Param("type")
	resourceID := c.Param("id")

	item, err := dao.GetResourceHealth(db, resourceType, resourceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Resource health not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: mergeResourceHealth(resourceType, resourceID, item.ResourceName, item)})
}

func ListExecutionLogs(c *gin.Context) {
	db := util.GetDatabase()
	query := db.Model(&dao.ExecutionLog{})
	if v := c.Query("resource_type"); v != "" {
		query = query.Where("resource_type = ?", v)
	}
	if v := c.Query("resource_id"); v != "" {
		query = query.Where("resource_id = ?", v)
	}
	if v := c.Query("trigger"); v != "" {
		query = query.Where("trigger = ?", v)
	}
	if v := c.Query("status"); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := c.Query("error_kind"); v != "" {
		query = query.Where("error_kind = ?", v)
	}
	items, err := dao.ListExecutionLogs(db, query.Limit(200))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: items})
}

func ResumeObservableResource(c *gin.Context) {
	resourceType := c.Param("type")
	resourceID := c.Param("id")
	if err := observability.ResumeResource(resourceType, resourceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Resource health not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Msg: "Success"})
}

func ListSystemNotifications(c *gin.Context) {
	db := util.GetDatabase()
	items, err := dao.ListSystemNotifications(db, db.Model(&dao.SystemNotification{}).Limit(200))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: items})
}

func SystemNotificationsRSS(c *gin.Context) {
	feed, err := observability.BuildNotificationFeed(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	rssFeed := feed.ToFeedsFeed()
	rssFeed.Link = &feeds.Link{Href: "/system/notifications/rss"}
	rss, err := rssFeed.ToRss()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rss))
}

func applyResourceHealthFilters(db *gorm.DB, resourceType string, status string) *gorm.DB {
	query := db.Model(&dao.ResourceHealth{})
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	if status != "" {
		query = query.Where("current_status = ?", status)
	}
	return query
}

func mergeResourceHealth(resourceType string, resourceID string, fallbackName string, item *dao.ResourceHealth) ResourceHealthView {
	view := ResourceHealthView{
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		ResourceName:  fallbackName,
		CurrentStatus: dao.ResourceStatusHealthy,
	}
	if item == nil {
		return view
	}
	if item.ResourceName != "" {
		view.ResourceName = item.ResourceName
	}
	view.CurrentStatus = item.CurrentStatus
	view.ConsecutiveFailures = item.ConsecutiveFailures
	view.LastSuccessAt = item.LastSuccessAt
	view.LastFailureAt = item.LastFailureAt
	view.LastErrorKind = item.LastErrorKind
	view.LastErrorMessage = item.LastErrorMessage
	view.PausedAt = item.PausedAt
	view.PausedReason = item.PausedReason
	return view
}
