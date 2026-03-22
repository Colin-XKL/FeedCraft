package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/feedruntime"
	"FeedCraft/internal/model"
	"FeedCraft/internal/observability"
	"FeedCraft/internal/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

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
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
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
