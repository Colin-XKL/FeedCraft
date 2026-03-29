package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/feedruntime"
	"FeedCraft/internal/model"
	"FeedCraft/internal/observability"
	"FeedCraft/internal/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

const topicDetailLogLimit = 20

var errTopicValidationRollback = errors.New("topic validation rollback")

type TopicValidationResult struct {
	Valid    bool                   `json:"valid"`
	Errors   []TopicValidationIssue `json:"errors"`
	Warnings []TopicValidationIssue `json:"warnings"`
}

type TopicValidationIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type TopicDetailResponse struct {
	Topic                dao.TopicFeed             `json:"topic"`
	PublicURL            string                    `json:"public_url"`
	Health               ResourceHealthView        `json:"health"`
	RecentExecutions     []TopicExecutionLogView   `json:"recent_executions"`
	RelatedNotifications []*dao.SystemNotification `json:"related_notifications"`
}

type TopicExecutionLogView struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	ResourceName string    `json:"resource_name"`
	Trigger      string    `json:"trigger"`
	Status       string    `json:"status"`
	ErrorKind    string    `json:"error_kind"`
	Message      string    `json:"message"`
	DetailsJSON  string    `json:"details_json"`
	Details      any       `json:"details,omitempty"`
	RequestID    string    `json:"request_id"`
	DurationMS   int64     `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}

func CreateTopicFeed(c *gin.Context) {
	var topicData dao.TopicFeed
	if err := c.ShouldBindJSON(&topicData); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateTopicFeed(db, &topicData); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: topicData})
}

func GetTopicFeed(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	topicData, err := dao.GetTopicFeedByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: topicData})
}

func ListTopicFeeds(c *gin.Context) {
	db := util.GetDatabase()
	topicList, err := dao.ListTopicFeeds(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: topicList})
}

func UpdateTopicFeed(c *gin.Context) {
	id := c.Param("id")
	var topicData dao.TopicFeed
	if err := c.ShouldBindJSON(&topicData); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// Ensure the ID in the URL matches the ID in the body
	if id != topicData.ID {
		topicData.ID = id
	}

	db := util.GetDatabase()

	_, err := dao.GetTopicFeedByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := dao.UpdateTopicFeed(db, &topicData); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: topicData})
}

func DeleteTopicFeed(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	if err := dao.DeleteTopicFeed(db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
}

func ValidateTopicFeed(c *gin.Context) {
	var topicData dao.TopicFeed
	if err := c.ShouldBindJSON(&topicData); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	result, err := validateTopicConfig(c.Request.Context(), util.GetDatabase(), &topicData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: result})
}

func GetTopicFeedDetail(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	topicData, err := dao.GetTopicFeedByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	health, err := dao.GetResourceHealth(db, dao.ResourceTypeTopic, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	executions, err := dao.ListExecutionLogs(
		db,
		db.Model(&dao.ExecutionLog{}).
			Where("resource_type = ? AND resource_id = ?", dao.ResourceTypeTopic, id).
			Limit(topicDetailLogLimit),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	notifications, err := dao.ListSystemNotifications(
		db,
		db.Model(&dao.SystemNotification{}).
			Where("resource_type = ? AND resource_id = ?", dao.ResourceTypeTopic, id).
			Limit(topicDetailLogLimit),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	detail := TopicDetailResponse{
		Topic:                *topicData,
		PublicURL:            "/topic/" + topicData.ID,
		Health:               mergeResourceHealth(dao.ResourceTypeTopic, topicData.ID, topicData.Title, health),
		RecentExecutions:     buildTopicExecutionViews(executions),
		RelatedNotifications: notifications,
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: detail})
}

func PublicTopicFeed(c *gin.Context) {
	topicID := c.Param("id")
	requestID := fmt.Sprintf("topic-%d", time.Now().UnixNano())
	ctx := observability.WithRequestID(c.Request.Context(), requestID)

	provider, err := feedruntime.NewBuilder(util.GetDatabase()).BuildTopicProvider(ctx, topicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	feed, err := provider.Fetch(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Topic feed not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	renderCraftFeedAsRSS(c, feed)
}

func renderCraftFeedAsRSS(c *gin.Context, feed *model.CraftFeed) {
	rssFeed := feed.ToFeedsFeed()
	rssStr, err := rssFeed.ToRss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to render RSS: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rssStr))
}

func validateTopicConfig(ctx context.Context, db *gorm.DB, topicData *dao.TopicFeed) (*TopicValidationResult, error) {
	result := &TopicValidationResult{
		Valid:    true,
		Errors:   []TopicValidationIssue{},
		Warnings: []TopicValidationIssue{},
	}
	if topicData == nil {
		result.Valid = false
		result.Errors = append(result.Errors, TopicValidationIssue{
			Field:   "topic",
			Message: "Topic config is required",
		})
		return result, nil
	}

	if strings.TrimSpace(topicData.ID) == "" {
		result.Valid = false
		result.Errors = append(result.Errors, TopicValidationIssue{
			Field:   "id",
			Message: "Topic ID is required",
		})
		return result, nil
	}

	if len(topicData.InputURIs) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, TopicValidationIssue{
			Field:   "input_uris",
			Message: "At least one input source is required",
		})
	}

	for idx, uri := range topicData.InputURIs {
		if strings.TrimSpace(uri) == "" {
			result.Valid = false
			result.Errors = append(result.Errors, TopicValidationIssue{
				Field:   fmt.Sprintf("input_uris[%d]", idx),
				Message: "Input URI cannot be empty",
			})
		}
	}

	if !result.Valid {
		return result, nil
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if saveErr := tx.Save(topicData).Error; saveErr != nil {
			return saveErr
		}

		_, buildErr := feedruntime.NewBuilder(tx).BuildTopicProvider(ctx, topicData.ID)
		if buildErr != nil {
			result.Valid = false
			result.Errors = append(result.Errors, TopicValidationIssue{
				Field:   "runtime",
				Message: buildErr.Error(),
			})
		}

		return errTopicValidationRollback
	})
	if err != nil && !errors.Is(err, errTopicValidationRollback) {
		return nil, err
	}

	return result, nil
}

func buildTopicExecutionViews(items []*dao.ExecutionLog) []TopicExecutionLogView {
	result := make([]TopicExecutionLogView, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		view := TopicExecutionLogView{
			ID:           item.ID,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			ResourceName: item.ResourceName,
			Trigger:      item.Trigger,
			Status:       item.Status,
			ErrorKind:    item.ErrorKind,
			Message:      item.Message,
			DetailsJSON:  item.DetailsJSON,
			RequestID:    item.RequestID,
			DurationMS:   item.DurationMS,
			CreatedAt:    item.CreatedAt,
		}
		if strings.TrimSpace(item.DetailsJSON) != "" {
			var parsed any
			if err := json.Unmarshal([]byte(item.DetailsJSON), &parsed); err == nil {
				view.Details = parsed
			}
		}
		result = append(result, view)
	}
	return result
}
