package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		Content: "before\u0005after",
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

	var count int64
	require.NoError(t, db.Model(&dao.InboxItem{}).Count(&count).Error)
	assert.Zero(t, count)
}

func inboxPushTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	util.SetDatabaseForTest(db)
	return db
}
