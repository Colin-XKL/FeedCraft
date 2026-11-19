package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const inboxPushShortContentThresholdBytes = 800

type InboxPushItem struct {
	Title     string  `json:"title"`
	URL       string  `json:"url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Summary   string  `json:"summary,omitempty"`
	ID        string  `json:"id,omitempty"`
	Author    string  `json:"author,omitempty"`
	Timestamp float64 `json:"timestamp,omitempty"`
	IsNotice  bool    `json:"is_notice,omitempty"`
}

type InboxPushItemResult struct {
	Index        int    `json:"index"`
	ID           string `json:"id,omitempty"`
	Status       string `json:"status"`
	ContentBytes int    `json:"content_bytes"`
	Error        string `json:"error,omitempty"`
	Warning      string `json:"warning,omitempty"`
}

type InboxPushResponse struct {
	Total   int                   `json:"total"`
	Created int                   `json:"created"`
	Updated int                   `json:"updated"`
	Failed  int                   `json:"failed"`
	Results []InboxPushItemResult `json:"results,omitempty"`
}

func PushInboxItems(c *gin.Context) {
	inboxID := c.Param("inbox_id")

	skipInvalid, err := parseBoolQueryParam(c, "skip_invalid")
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}
	requireContent, err := parseBoolQueryParam(c, "require_content")
	if err != nil {
		c.JSON(http.StatusBadRequest, util.APIResponse[any]{Msg: err.Error()})
		return
	}

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

	results := make([]InboxPushItemResult, len(reqItems))
	validIndices := make([]int, 0, len(reqItems))
	failedCount := 0
	firstErrMsg := ""

	for i, item := range reqItems {
		trimmedID := strings.TrimSpace(item.ID)
		contentBytes := len(item.Content)
		res := InboxPushItemResult{
			Index:        i,
			ID:           trimmedID,
			Status:       "ok",
			ContentBytes: contentBytes,
		}

		if errMsg := validateInboxPushItem(i, item, requireContent); errMsg != "" {
			res.Status = "error"
			res.Error = errMsg
			failedCount++
			if firstErrMsg == "" {
				firstErrMsg = errMsg
			}
		} else {
			res.Warning = inboxPushContentWarning(item)
			validIndices = append(validIndices, i)
		}
		results[i] = res
	}

	if failedCount > 0 && (!skipInvalid || len(validIndices) == 0) {
		c.JSON(http.StatusBadRequest, util.APIResponse[InboxPushResponse]{
			Msg: firstErrMsg,
			Data: InboxPushResponse{
				Total:   len(reqItems),
				Created: 0,
				Updated: 0,
				Failed:  failedCount,
				Results: results,
			},
		})
		return
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

	// Extract explicitly provided IDs from valid items to count updates
	var explicitIDs []string
	for _, idx := range validIndices {
		if trimmedID := strings.TrimSpace(reqItems[idx].ID); trimmedID != "" {
			explicitIDs = append(explicitIDs, trimmedID)
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

	daoItems := make([]dao.InboxItem, 0, len(validIndices))
	siteBaseURL := strings.TrimRight(util.GetEnvClient().GetString("SITE_BASE_URL"), "/")
	now := time.Now()

	updatedCount := 0
	for _, idx := range validIndices {
		req := reqItems[idx]
		itemID := strings.TrimSpace(req.ID)
		if itemID == "" {
			itemID = uuid.New().String()
		} else if existingMap[itemID] {
			updatedCount++
		}
		results[idx].ID = itemID

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
		Created: len(validIndices) - updatedCount,
		Updated: updatedCount,
		Failed:  failedCount,
		Results: results,
	})
}

func parseBoolQueryParam(c *gin.Context, name string) (bool, error) {
	raw := strings.TrimSpace(c.Query(name))
	if raw == "" {
		return false, nil
	}
	val, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("invalid boolean query parameter %q: %q", name, raw)
	}
	return val, nil
}

func validateInboxPushItem(index int, item InboxPushItem, requireContent bool) string {
	if strings.TrimSpace(item.Title) == "" {
		return fmt.Sprintf("Item at index %d is missing required 'title'", index)
	}
	if field, invalidChar, characterPosition, byteOffset, ok := inboxPushInvalidXMLChar(item); ok {
		return fmt.Sprintf(
			"Item at index %d contains invalid XML 1.0 character U+%04X in field %q at character %d (UTF-8 byte offset %d)",
			index,
			invalidChar,
			field,
			characterPosition,
			byteOffset,
		)
	}
	if requireContent && !item.IsNotice && strings.TrimSpace(item.Content) == "" {
		return fmt.Sprintf("Item at index %d is missing required 'content'", index)
	}
	return ""
}

func inboxPushContentWarning(item InboxPushItem) string {
	if item.IsNotice {
		return ""
	}
	contentBytes := len(item.Content)
	if contentBytes == 0 {
		return "content is empty (0 bytes); item may fall back to summary only (set is_notice=true for notice-only items)"
	}
	if strings.TrimSpace(item.Content) == "" {
		return fmt.Sprintf("content is blank (%d bytes); item may fall back to summary only (set is_notice=true for notice-only items)", contentBytes)
	}
	if contentBytes < inboxPushShortContentThresholdBytes {
		return fmt.Sprintf(
			"content is short (%d bytes, below %d-byte threshold); verify full content was sent (set is_notice=true for notice-only items)",
			contentBytes,
			inboxPushShortContentThresholdBytes,
		)
	}
	return ""
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
