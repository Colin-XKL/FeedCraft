package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPushInboxItemsRejectsXMLInvalidControlCharacters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "xml-validation-inbox",
		MaxItems: 10,
	}))

	body, err := json.Marshal([]InboxPushItem{{
		ID:      "invalid-control-character",
		Title:   "Invalid XML",
		Content: "前缀\u0005after",
	}})
	require.NoError(t, err)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/inbox/xml-validation-inbox/items",
		bytes.NewReader(body),
	)
	recorder := httptest.NewRecorder()
	router := gin.New()
	router.POST("/api/inbox/:inbox_id/items", PushInboxItems)
	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "content")
	assert.Contains(t, recorder.Body.String(), "U+0005")
	assert.Contains(t, recorder.Body.String(), "character 3")
	assert.Contains(t, recorder.Body.String(), "UTF-8 byte offset 6")

	var resp util.APIResponse[InboxPushResponse]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Data.Total)
	assert.Equal(t, 1, resp.Data.Failed)
	require.Len(t, resp.Data.Results, 1)
	assert.Equal(t, "error", resp.Data.Results[0].Status)
	assert.Equal(t, "invalid-control-character", resp.Data.Results[0].ID)

	var count int64
	require.NoError(t, db.Model(&dao.InboxItem{}).Count(&count).Error)
	assert.Zero(t, count)
}

func TestPushInboxItemsBatchValidationReportsAllItemsWithoutSkipInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "batch-validation-inbox",
		MaxItems: 10,
	}))

	longContent := strings.Repeat("a", 900)
	body, err := json.Marshal([]InboxPushItem{
		{ID: "good-1", Title: "Valid Item", Content: longContent},
		{ID: "bad-xml", Title: "Bad XML", Content: "abc\u0000def"},
		{ID: "bad-title", Title: "   ", Content: longContent},
	})
	require.NoError(t, err)

	recorder := performInboxPushRequest("/api/inbox/batch-validation-inbox/items", body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp util.APIResponse[InboxPushResponse]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Contains(t, resp.Msg, "U+0000")
	assert.Equal(t, 3, resp.Data.Total)
	assert.Equal(t, 0, resp.Data.Created)
	assert.Equal(t, 0, resp.Data.Updated)
	assert.Equal(t, 2, resp.Data.Failed)
	require.Len(t, resp.Data.Results, 3)

	assert.Equal(t, "ok", resp.Data.Results[0].Status)
	assert.Equal(t, 900, resp.Data.Results[0].ContentBytes)
	assert.Empty(t, resp.Data.Results[0].Error)
	assert.Empty(t, resp.Data.Results[0].Warning)

	assert.Equal(t, "error", resp.Data.Results[1].Status)
	assert.Contains(t, resp.Data.Results[1].Error, `field "content"`)
	assert.Contains(t, resp.Data.Results[1].Error, "U+0000")
	assert.Contains(t, resp.Data.Results[1].Error, "character 4")
	assert.Contains(t, resp.Data.Results[1].Error, "UTF-8 byte offset 3")

	assert.Equal(t, "error", resp.Data.Results[2].Status)
	assert.Contains(t, resp.Data.Results[2].Error, "missing required 'title'")

	var count int64
	require.NoError(t, db.Model(&dao.InboxItem{}).Count(&count).Error)
	assert.Zero(t, count)
}

func TestPushInboxItemsSkipInvalidPersistsValidItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "skip-invalid-inbox",
		MaxItems: 10,
	}))
	require.NoError(t, db.Create(&dao.InboxItem{
		InboxID: "skip-invalid-inbox",
		ItemID:  "existing-1",
		Title:   "Old Title",
		Content: "Old Content",
	}).Error)

	longContent := strings.Repeat("b", 850)
	body, err := json.Marshal([]InboxPushItem{
		{Title: "New Valid", Content: longContent},
		{ID: "bad-xml", Title: "Bad XML", Content: "前缀\u0005after"},
		{ID: "existing-1", Title: "Updated Valid", Content: longContent},
	})
	require.NoError(t, err)

	recorder := performInboxPushRequest("/api/inbox/skip-invalid-inbox/items?skip_invalid=true", body)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var resp InboxPushResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Equal(t, 3, resp.Total)
	assert.Equal(t, 1, resp.Created)
	assert.Equal(t, 1, resp.Updated)
	assert.Equal(t, 1, resp.Failed)
	require.Len(t, resp.Results, 3)

	assert.Equal(t, "ok", resp.Results[0].Status)
	assert.NotEmpty(t, resp.Results[0].ID)
	assert.Equal(t, 850, resp.Results[0].ContentBytes)

	assert.Equal(t, "error", resp.Results[1].Status)
	assert.Equal(t, "bad-xml", resp.Results[1].ID)
	assert.Contains(t, resp.Results[1].Error, "U+0005")
	assert.Contains(t, resp.Results[1].Error, "character 3")
	assert.Contains(t, resp.Results[1].Error, "UTF-8 byte offset 6")

	assert.Equal(t, "ok", resp.Results[2].Status)
	assert.Equal(t, "existing-1", resp.Results[2].ID)

	items, err := dao.ListInboxItems(db, "skip-invalid-inbox")
	require.NoError(t, err)
	assert.Len(t, items, 2)

	updatedItem, err := dao.GetInboxItemByItemID(db, "skip-invalid-inbox", "existing-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated Valid", updatedItem.Title)
}

func TestPushInboxItemsSkipInvalidAllInvalidReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "all-invalid-inbox",
		MaxItems: 10,
	}))

	body, err := json.Marshal([]InboxPushItem{
		{ID: "bad-1", Title: "", Content: "hello"},
		{ID: "bad-2", Title: "Bad", Content: "x\u0007y"},
	})
	require.NoError(t, err)

	recorder := performInboxPushRequest("/api/inbox/all-invalid-inbox/items?skip_invalid=true", body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp util.APIResponse[InboxPushResponse]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Equal(t, 2, resp.Data.Total)
	assert.Equal(t, 2, resp.Data.Failed)
	require.Len(t, resp.Data.Results, 2)
	assert.Equal(t, "error", resp.Data.Results[0].Status)
	assert.Equal(t, "error", resp.Data.Results[1].Status)
}

func TestPushInboxItemsContentBytesAndShortContentWarning(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "content-warning-inbox",
		MaxItems: 10,
	}))

	fullContent := strings.Repeat("x", inboxPushShortContentThresholdBytes)
	shortContent := strings.Repeat("y", 120)

	body, err := json.Marshal([]InboxPushItem{
		{ID: "full", Title: "Full Report", Content: fullContent},
		{ID: "missing", Title: "Missing Content", Summary: "Only summary"},
		{ID: "short", Title: "Short Content", Content: shortContent},
		{ID: "notice", Title: "Pure Notice", Summary: "System alert", IsNotice: true},
	})
	require.NoError(t, err)

	recorder := performInboxPushRequest("/api/inbox/content-warning-inbox/items", body)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var resp InboxPushResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	assert.Equal(t, 4, resp.Total)
	assert.Equal(t, 4, resp.Created)
	assert.Equal(t, 0, resp.Failed)
	require.Len(t, resp.Results, 4)

	assert.Equal(t, inboxPushShortContentThresholdBytes, resp.Results[0].ContentBytes)
	assert.Empty(t, resp.Results[0].Warning)

	assert.Equal(t, 0, resp.Results[1].ContentBytes)
	assert.Contains(t, resp.Results[1].Warning, "empty")

	assert.Equal(t, 120, resp.Results[2].ContentBytes)
	assert.Contains(t, resp.Results[2].Warning, "short")

	assert.Equal(t, 0, resp.Results[3].ContentBytes)
	assert.Empty(t, resp.Results[3].Warning)
}

func TestPushInboxItemsRequireContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := inboxPushTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.Inbox{}, &dao.InboxItem{}))
	require.NoError(t, dao.CreateInbox(db, &dao.Inbox{
		ID:       "require-content-inbox",
		MaxItems: 10,
	}))

	fullContent := strings.Repeat("z", 900)

	t.Run("rejects missing content by default", func(t *testing.T) {
		body, err := json.Marshal([]InboxPushItem{
			{ID: "with-content", Title: "Has Content", Content: fullContent},
			{ID: "no-content", Title: "No Content", Summary: "Summary only"},
		})
		require.NoError(t, err)

		recorder := performInboxPushRequest("/api/inbox/require-content-inbox/items?require_content=true", body)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)

		var resp util.APIResponse[InboxPushResponse]
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
		assert.Equal(t, 1, resp.Data.Failed)
		assert.Equal(t, "ok", resp.Data.Results[0].Status)
		assert.Equal(t, "error", resp.Data.Results[1].Status)
		assert.Contains(t, resp.Data.Results[1].Error, "missing required 'content'")
	})

	t.Run("allows explicit notice items even when require_content is true", func(t *testing.T) {
		body, err := json.Marshal([]InboxPushItem{
			{ID: "notice-only", Title: "Alert", Summary: "Alert body", IsNotice: true},
		})
		require.NoError(t, err)

		recorder := performInboxPushRequest("/api/inbox/require-content-inbox/items?require_content=true", body)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		var resp InboxPushResponse
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
		assert.Equal(t, 1, resp.Created)
		assert.Empty(t, resp.Results[0].Warning)
	})

	t.Run("combines require_content and skip_invalid", func(t *testing.T) {
		body, err := json.Marshal([]InboxPushItem{
			{ID: "good-content", Title: "Good", Content: fullContent},
			{ID: "missing-content", Title: "Missing", Summary: "Only summary"},
		})
		require.NoError(t, err)

		recorder := performInboxPushRequest("/api/inbox/require-content-inbox/items?require_content=true&skip_invalid=true", body)
		assert.Equal(t, http.StatusCreated, recorder.Code)

		var resp InboxPushResponse
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
		assert.Equal(t, 2, resp.Total)
		assert.Equal(t, 1, resp.Created)
		assert.Equal(t, 1, resp.Failed)
		assert.Equal(t, "ok", resp.Results[0].Status)
		assert.Equal(t, "error", resp.Results[1].Status)
	})
}

func TestPushInboxItemsInvalidQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body, err := json.Marshal([]InboxPushItem{{Title: "Test", Content: "body"}})
	require.NoError(t, err)

	recorder := performInboxPushRequest("/api/inbox/any/items?skip_invalid=notabool", body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "skip_invalid")

	recorder = performInboxPushRequest("/api/inbox/any/items?require_content=notabool", body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "require_content")
}

func performInboxPushRequest(target string, body []byte) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	router := gin.New()
	router.POST("/api/inbox/:inbox_id/items", PushInboxItems)
	router.ServeHTTP(recorder, request)
	return recorder
}

func inboxPushTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	util.SetDatabaseForTest(db)
	return db
}
