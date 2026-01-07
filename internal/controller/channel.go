package controller

import (
	"FeedCraft/internal/dao"
	"FeedCraft/internal/recipe"
	"FeedCraft/internal/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func CreateChannel(c *gin.Context) {
	var channel dao.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateChannel(db, &channel); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: channel})
}

func GetChannel(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	channel, err := dao.GetChannelByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: channel})
}

// ChannelInfo detail info including dao.Channel basic info and preheat stats
type ChannelInfo struct {
	ID             string    `json:"id,omitempty" binding:"required"`
	Description    string    `json:"description,omitempty"`
	ProcessorName  string    `json:"processor_name" binding:"required"` // Formerly Craft
	SourceType     string    `json:"source_type"`
	SourceConfig   string    `json:"source_config"`
	IsActive       bool      `json:"is_active" `
	LastAccessedAt time.Time `json:"last_accessed_at"`
}

func ListChannels(c *gin.Context) {
	db := util.GetDatabase()
	channels, err := dao.ListChannels(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	channelInfoList := lo.Map(channels, func(item *dao.Channel, index int) ChannelInfo {
		recipeStatus := recipe.Scheduler.GetContextInfo(item.ID)
		return ChannelInfo{
			ID:             item.ID,
			Description:    item.Description,
			ProcessorName:  item.ProcessorName,
			SourceType:     item.SourceType,
			SourceConfig:   item.SourceConfig,
			IsActive:       recipeStatus.IsActive,
			LastAccessedAt: recipeStatus.LastRequestTime,
		}
	})
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: channelInfoList})
}

func UpdateChannel(c *gin.Context) {
	id := c.Param("id")
	var channel dao.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	_, err := dao.GetChannelByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Channel not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := dao.UpdateChannel(db, &channel); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: channel})
}

func DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	if err := dao.DeleteChannel(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{})
}
