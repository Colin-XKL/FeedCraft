package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InboxPushItem struct {
	Title     string  `json:"title"`
	URL       string  `json:"url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Summary   string  `json:"summary,omitempty"`
	ID        string  `json:"id,omitempty"`
	Author    string  `json:"author,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
}

type InboxPushResponse struct {
	Total   int `json:"total"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

func PushInboxItems(c *gin.Context) {
	inboxID := c.Param("inbox_id")
	var reqItems []InboxPushItem

	if err := c.ShouldBindJSON(&reqItems); err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Invalid JSON array: " + err.Error()})
		return
	}

	if len(reqItems) == 0 {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Array cannot be empty"})
		return
	}

	if len(reqItems) > 100 {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: "Batch size exceeds maximum limit of 100"})
		return
	}

	for i, item := range reqItems {
		if strings.TrimSpace(item.Title) == "" {
			c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: fmt.Sprintf("Item at index %d is missing required 'title'", i)})
			return
		}
		if field, invalidChar, characterPosition, byteOffset, ok := inboxPushInvalidXMLChar(item); ok {
			c.JSON(http.StatusBadRequest, util.APIResponse[any]{
				Msg: fmt.Sprintf(
					"Item at index %d contains invalid XML 1.0 character U+%04X in field %q at character %d (UTF-8 byte offset %d)",
					i,
					invalidChar,
					field,
					characterPosition,
					byteOffset,
				),
			})
			return
		}
	}

	db := util.GetDatabase()

	inbox, err := dao.GetInboxByID(db, inboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, util.APIResponse[any]{Msg: "Inbox not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: err.Error()})
		return
	}

	// Extract explicitly provided IDs to count updates
	var explicitIDs []string
	for _, item := range reqItems {
		if strings.TrimSpace(item.ID) != "" {
			explicitIDs = append(explicitIDs, strings.TrimSpace(item.ID))
		}
	}

	existingIDs, err := dao.ListExistingInboxItemIDs(db, inboxID, explicitIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Database error checking existing items"})
		return
	}
	existingMap := make(map[string]bool)
	for _, id := range existingIDs {
		existingMap[id] = true
	}

	var daoItems []dao.InboxItem
	siteBaseURL := strings.TrimRight(util.GetEnvClient().GetString("SITE_BASE_URL"), "/")
	now := time.Now()

	updatedCount := 0
	for _, req := range reqItems {
		itemID := strings.TrimSpace(req.ID)
		if itemID == "" {
			itemID = uuid.New().String()
		} else if existingMap[itemID] {
			updatedCount++
		}

		pubTime := now
		if req.Timestamp > 0 {
			pubTime = time.Unix(int64(req.Timestamp), 0)
		}

		summary := strings.TrimSpace(req.Summary)
		if summary == "" {
			runes := []rune(req.Content)
			if len(runes) > 200 {
				summary = string(runes[:200])
			} else {
				summary = req.Content
			}
		}

		itemURL := strings.TrimSpace(req.URL)
		if itemURL == "" {
			itemURL = fmt.Sprintf("%s/inbox/%s/items/%s/content", siteBaseURL, inboxID, itemID)
		}

		daoItems = append(daoItems, dao.InboxItem{
			InboxID:     inboxID,
			ItemID:      itemID,
			Title:       req.Title,
			URL:         itemURL,
			Content:     req.Content,
			Summary:     summary,
			Author:      req.Author,
			PublishedAt: pubTime,
			CreatedAt:   now,
		})
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// Bulk Upsert
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "inbox_id"}, {Name: "item_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"title", "url", "content", "summary", "author", "published_at",
			}),
		}).Create(&daoItems).Error; err != nil {
			return err
		}

		// Enforce MaxItems
		if err := dao.DeleteOverflowInboxItems(tx, inboxID, inbox.MaxItems); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.APIResponse[any]{Msg: "Failed to persist items: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, InboxPushResponse{
		Total:   len(reqItems),
		Created: len(reqItems) - updatedCount,
		Updated: updatedCount,
	})
}

func inboxPushInvalidXMLChar(item InboxPushItem) (string, rune, int, int, bool) {
	fields := []struct {
		name  string
		value string
	}{
		{name: "title", value: item.Title},
		{name: "url", value: item.URL},
		{name: "content", value: item.Content},
		{name: "summary", value: item.Summary},
		{name: "id", value: item.ID},
		{name: "author", value: item.Author},
	}
	for _, field := range fields {
		characterPosition := 1
		for byteOffset, char := range field.value {
			if !isXML10Char(char) {
				return field.name, char, characterPosition, byteOffset, true
			}
			characterPosition++
		}
	}
	return "", 0, 0, 0, false
}

func isXML10Char(char rune) bool {
	return char == '\t' || char == '\n' || char == '\r' ||
		(char >= 0x20 && char <= 0xD7FF) ||
		(char >= 0xE000 && char <= 0xFFFD) ||
		(char >= 0x10000 && char <= 0x10FFFF)
}
