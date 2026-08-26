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

func TestCreateCustomRecipe_DuplicateIDReturnsConflictWithoutSQL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := customRecipeTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}))

	router := gin.New()
	router.POST("/api/admin/recipes", CreateCustomRecipe)

	payload := map[string]string{
		"id":            "e2e-ai-recipe-887",
		"description":   "first recipe",
		"craft":         "fulltext",
		"source_type":   "rss",
		"source_config": `{"http_fetcher":{"url":"https://example.com/rss.xml"}}`,
	}

	first := postCustomRecipe(t, router, payload)
	require.Equal(t, http.StatusCreated, first.Code, first.Body.String())

	payload["description"] = "duplicate recipe"
	second := postCustomRecipe(t, router, payload)

	assert.Equal(t, http.StatusConflict, second.Code, second.Body.String())

	var response util.APIResponse[any]
	require.NoError(t, json.Unmarshal(second.Body.Bytes(), &response))
	assert.NotEmpty(t, response.Msg)
	assert.NotContains(t, strings.ToLower(response.Msg), "unique constraint")
	assert.NotContains(t, strings.ToLower(response.Msg), "sqlite")
	assert.NotContains(t, response.Msg, "custom_recipes_v2")
	assert.Contains(t, strings.ToLower(response.Msg), "already exists")
}

func TestCreateCustomRecipe_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := customRecipeTestDatabase(t)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}))

	router := gin.New()
	router.POST("/api/admin/recipes", CreateCustomRecipe)

	recorder := postCustomRecipe(t, router, map[string]string{
		"id":            "Bad Name",
		"craft":         "fulltext",
		"source_type":   "rss",
		"source_config": `{"http_fetcher":{"url":"https://example.com/rss.xml"}}`,
	})

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assertJSONMessageContains(t, recorder, "lowercase")
}

func customRecipeTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	util.SetDatabaseForTest(db)
	return db
}

func postCustomRecipe(t *testing.T, router *gin.Engine, payload map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/recipes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}
