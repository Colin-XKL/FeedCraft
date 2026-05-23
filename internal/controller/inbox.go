package controller

import (
	"errors"
	"net/http"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	if err := dao.DeleteOverflowInboxItems(db, id, data.MaxItems); err != nil {
		logrus.Warnf("failed to prune overflow items on update for inbox %s: %v", id, err)
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{Data: data})
}

func GetInboxGCStats(c *gin.Context) {
	db := util.GetDatabase()

	var totalItems int64
	if err := db.Model(&dao.InboxItem{}).Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	var orphanedCount int64
	if err := db.Model(&dao.InboxItem{}).
		Where("inbox_id NOT IN (SELECT id FROM inboxes)").
		Count(&orphanedCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	var inboxes []dao.Inbox
	if err := db.Find(&inboxes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	var overflowCount int64
	for _, inbox := range inboxes {
		var count int64
		if err := db.Model(&dao.InboxItem{}).Where("inbox_id = ?", inbox.ID).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
			return
		}
		if count > int64(inbox.MaxItems) {
			overflowCount += (count - int64(inbox.MaxItems))
		}
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{
		Data: gin.H{
			"total_items":    totalItems,
			"orphaned_count": orphanedCount,
			"overflow_count": overflowCount,
		},
	})
}

func TriggerInboxGCCleanup(c *gin.Context) {
	db := util.GetDatabase()

	var orphanedDeleted int64
	var overflowDeleted int64

	err := db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("inbox_id NOT IN (SELECT id FROM inboxes)").Delete(&dao.InboxItem{})
		if res.Error != nil {
			return res.Error
		}
		orphanedDeleted = res.RowsAffected

		var inboxes []dao.Inbox
		if err := tx.Find(&inboxes).Error; err != nil {
			return err
		}

		for _, inbox := range inboxes {
			var ids []uint
			if err := tx.Model(&dao.InboxItem{}).
				Where("inbox_id = ?", inbox.ID).
				Order("created_at desc").
				Offset(inbox.MaxItems).
				Pluck("id", &ids).Error; err != nil {
				return err
			}
			if len(ids) > 0 {
				resDel := tx.Where("id IN ?", ids).Delete(&dao.InboxItem{})
				if resDel.Error != nil {
					return resDel.Error
				}
				overflowDeleted += resDel.RowsAffected
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, util.APIResponse[any]{
		Data: gin.H{
			"orphaned_deleted": orphanedDeleted,
			"overflow_deleted": overflowDeleted,
		},
	})
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
