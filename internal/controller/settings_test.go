package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
)

func TestSaveSearchProviderConfig(t *testing.T) {
	// Setup DB
	tmpDir, err := os.MkdirTemp("", "feedcraft_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("DB_SQLITE_PATH", tmpDir)
	os.Setenv("FC_DB_SQLITE_PATH", tmpDir)
	db := util.GetDatabase()
	if err := db.AutoMigrate(&dao.SystemSetting{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/save", SaveSearchProviderConfig)

	// Helper to make request
	makeRequest := func(body interface{}) *httptest.ResponseRecorder {
		jsonBytes, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/save", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	// 1. Initial Save with API Key
	initialConfig := SearchProviderConfigRequest{
		SearchProviderConfig: config.SearchProviderConfig{
			Provider: "litellm",
			APIKey:   "initial-key",
			APIUrl:   "http://example.com",
		},
		UpdateAPIKey: true,
	}
	w := makeRequest(initialConfig)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify in DB
	var savedCfg config.SearchProviderConfig
	err = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if err != nil {
		t.Fatalf("Failed to get setting: %v", err)
	}
	if savedCfg.APIKey != "initial-key" {
		t.Errorf("Expected initial-key, got %s", savedCfg.APIKey)
	}

	// 2. Save with empty key and UpdateAPIKey=false (Legacy/Keep)
	keepConfig := SearchProviderConfigRequest{
		SearchProviderConfig: config.SearchProviderConfig{
			Provider: "litellm",
			APIKey:   "", // Empty
			APIUrl:   "http://example.com",
		},
		UpdateAPIKey: false,
	}
	w = makeRequest(keepConfig)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if savedCfg.APIKey != "initial-key" {
		t.Errorf("Expected key to persist (initial-key), got %s", savedCfg.APIKey)
	}

	// 3. Save with empty key and UpdateAPIKey=true (Clear)
	clearConfig := SearchProviderConfigRequest{
		SearchProviderConfig: config.SearchProviderConfig{
			Provider: "litellm",
			APIKey:   "", // Empty
			APIUrl:   "http://example.com",
		},
		UpdateAPIKey: true,
	}
	w = makeRequest(clearConfig)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if savedCfg.APIKey != "" {
		t.Errorf("Expected key to be cleared, got %s", savedCfg.APIKey)
	}

	// 4. Save with new key and UpdateAPIKey=true
	newKeyConfig := SearchProviderConfigRequest{
		SearchProviderConfig: config.SearchProviderConfig{
			Provider: "litellm",
			APIKey:   "new-key",
			APIUrl:   "http://example.com",
		},
		UpdateAPIKey: true,
	}
	w = makeRequest(newKeyConfig)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if savedCfg.APIKey != "new-key" {
		t.Errorf("Expected new-key, got %s", savedCfg.APIKey)
	}

	// 5. Save with new key and UpdateAPIKey=false (Should update because APIKey is not empty)
	implicitUpdateConfig := SearchProviderConfigRequest{
		SearchProviderConfig: config.SearchProviderConfig{
			Provider: "litellm",
			APIKey:   "newer-key",
			APIUrl:   "http://example.com",
		},
		UpdateAPIKey: false,
	}
	w = makeRequest(implicitUpdateConfig)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
	dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if savedCfg.APIKey != "newer-key" {
		t.Errorf("Expected newer-key, got %s", savedCfg.APIKey)
	}
}
