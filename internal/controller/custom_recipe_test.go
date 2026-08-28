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
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCustomRecipeTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:custom_recipe_test?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&dao.CustomRecipeV2{}))
	util.SetDatabaseForTest(db)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/recipes", CreateCustomRecipe)
	return router
}

func TestCreateCustomRecipeRejectsEmptyRequiredFields(t *testing.T) {
	router := setupCustomRecipeTestRouter(t)

	body, err := json.Marshal(map[string]string{
		"id":            "",
		"craft":         "",
		"source_type":   "rss",
		"source_config": `{"http_fetcher":{"url":"https://example.com/rss.xml"}}`,
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var resp util.APIResponse[any]
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &resp))
	require.Contains(t, resp.Msg, "name is required")
	require.Contains(t, resp.Msg, "craft is required")
	require.NotContains(t, resp.Msg, "Field validation")
	require.NotContains(t, resp.Msg, "CustomRecipeV2")
}
