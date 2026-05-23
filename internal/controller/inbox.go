package controller

import (
	"errors"
	"net/http"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateInbox(c *gin.Context) {
	var data dao.Inbox
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	db := util.GetDatabase()

	if err := dao.CreateInbox(db, &data); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: data})
}

func GetInbox(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	data, err := dao.GetInboxByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Inbox not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: data})
}

func UpdateInbox(c *gin.Context) {
	id := c.Param("id")
	var data dao.Inbox
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	data.ID = id
	db := util.GetDatabase()

	if _, err := dao.GetInboxByID(db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Inbox not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	if err := dao.UpdateInbox(db, &data); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: data})
}

func DeleteInbox(c *gin.Context) {
	id := c.Param("id")
	db := util.GetDatabase()

	if err := dao.DeleteInbox(db, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Inbox not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{})
}

func ListInboxes(c *gin.Context) {
	db := util.GetDatabase()
	list, err := dao.ListInboxes(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: list})
}

// SystemAuthToken Admin CRUD

func CreateSystemAuthToken(c *gin.Context) {
	var data struct {
		Label string `json:"label"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	tokenValue := uuid.New().String()
	tokenRecord := dao.SystemAuthToken{
		Token: tokenValue,
		Label: data.Label,
	}

	db := util.GetDatabase()
	if err := dao.CreateSystemAuthToken(db, &tokenRecord); err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, util.APIResponse[any]{Data: tokenRecord})
}

func ListSystemAuthTokens(c *gin.Context) {
	db := util.GetDatabase()
	list, err := dao.ListSystemAuthTokens(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{Data: list})
}

func DeleteSystemAuthToken(c *gin.Context) {
	idParam := c.Param("id")
	db := util.GetDatabase()

	if err := db.Where("id = ?", idParam).Delete(&dao.SystemAuthToken{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, util.APIResponse[any]{})
}
