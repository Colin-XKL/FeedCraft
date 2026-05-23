package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/model"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// checkInboxAccess validates access for a potentially private inbox.
// It writes the HTTP error response and returns a non-nil error when access is denied.
// Callers should abort immediately on a non-nil return value.
func checkInboxAccess(c *gin.Context, inbox *dao.Inbox, db *gorm.DB) error {
	if inbox.IsPublic != nil && *inbox.IsPublic {
		return nil
	}

	tokenValue := c.Query("token")
	if tokenValue == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
				tokenValue = strings.TrimSpace(parts[1])
			}
		}
	}

	if tokenValue == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization required for private inbox"})
		return fmt.Errorf("unauthorized")
	}

	_, err := dao.GetSystemAuthTokenByToken(db, tokenValue)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Invalid token"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Database error"})
		}
		return err
	}
	return nil
}

// PublicInboxItemContent serves the raw HTML content of a single inbox article.
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

	if err := checkInboxAccess(c, inbox, db); err != nil {
		return
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

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(item.Content))
}

// PublicInboxRSSFeed serves a direct RSS feed for an inbox without requiring a Custom Recipe.
// Access control mirrors PublicInboxItemContent: public inboxes are open, private ones need a token.
func PublicInboxRSSFeed(c *gin.Context) {
	inboxID := c.Param("inbox_id")
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

	if err := checkInboxAccess(c, inbox, db); err != nil {
		return
	}

	items, err := dao.ListInboxItems(db, inboxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	siteBaseURL := strings.TrimRight(util.GetEnvClient().GetString("SITE_BASE_URL"), "/")
	inboxBaseURL := fmt.Sprintf("%s/inbox/%s", siteBaseURL, inboxID)

	craftFeed := &model.CraftFeed{
		Title:       inbox.Title,
		Description: inbox.Description,
		Id:          inboxBaseURL + "/rss",
		Link:        inboxBaseURL,
		Articles:    make([]*model.CraftArticle, 0, len(items)),
	}

	for _, item := range items {
		craftFeed.Articles = append(craftFeed.Articles, &model.CraftArticle{
			Title:       item.Title,
			Link:        item.URL,
			Content:     item.Content,
			Description: item.Summary,
			Id:          item.ItemID,
			AuthorName:  item.Author,
			Created:     item.PublishedAt,
			Updated:     item.PublishedAt,
		})
	}

	rssStr, err := craftFeed.ToFeedsFeed().ToRss()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to render RSS"})
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8", []byte(rssStr))
}
