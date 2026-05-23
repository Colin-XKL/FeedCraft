package controller

import (
	"errors"
	"net/http"
	"strings"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PublicInboxItemContent(c *gin.Context) {
	inboxID := c.Param("inbox_id")
	articleID := c.Param("article_id")
	db := util.GetDatabase()

	inbox, err := dao.GetInboxByID(db, inboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Inbox not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	// 权限校验逻辑
	if inbox.IsPublic == nil || !*inbox.IsPublic {
		var tokenValue string

		// 1. 优先尝试从 Query 提取 "?token=xxx"
		if qToken := c.Query("token"); qToken != "" {
			tokenValue = qToken
		} else {
			// 2. 其次尝试从 Authorization Header 提取 "Bearer xxx"
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenValue = strings.TrimSpace(parts[1])
				}
			}
		}

		if tokenValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization required for private inbox"})
			return
		}

		// 3. 查库校验
		_, err := dao.GetSystemAuthTokenByToken(db, tokenValue)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid token"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Database error"})
			return
		}
	}

	item, err := dao.GetInboxItemByItemID(db, inboxID, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	content := item.Content
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}
