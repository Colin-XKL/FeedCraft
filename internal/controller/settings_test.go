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
	"FeedCraft/internal/favicon"
	"FeedCraft/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestSaveSearchProviderConfig(t *testing.T) {
	// Setup DB
	tmpDir, err := os.MkdirTemp("", "feedcraft_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	t.Setenv("DB_SQLITE_PATH", tmpDir)
	t.Setenv("FC_DB_SQLITE_PATH", tmpDir)
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
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
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
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
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
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
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
	_ = dao.GetJsonSetting(db, constant.KeySearchProviderConfig, &savedCfg)
	if savedCfg.APIKey != "newer-key" {
		t.Errorf("Expected newer-key, got %s", savedCfg.APIKey)
	}
}

func TestFaviconProviderSettingsAPI(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:favicon-settings-controller?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	util.SetDatabaseForTest(db)
	if err := db.AutoMigrate(&dao.SystemSetting{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	previous := favicon.Settings()
	t.Cleanup(func() {
		if err := favicon.Replace(previous); err != nil {
			t.Fatalf("restore favicon settings: %v", err)
		}
	})
	if err := favicon.Replace(favicon.DefaultSettings()); err != nil {
		t.Fatalf("reset favicon settings: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/settings", GetFaviconProviderConfig)
	r.POST("/settings", SaveFaviconProviderConfig)
	r.POST("/preview", PreviewFaviconProviderConfig)

	settings := config.FaviconSettings{
		DefaultProviderID: "self_hosted",
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "self_hosted",
				Name:        "Self hosted",
				URLTemplate: "https://icons.example.test/favicon?domain={host}&size={size}",
				Enabled:     true,
			},
		},
	}
	saveResponse := performJSONRequest(t, r, http.MethodPost, "/settings", settings)
	if saveResponse.Code != http.StatusOK {
		t.Fatalf("save status = %d, body = %s", saveResponse.Code, saveResponse.Body.String())
	}

	_, activeProvider := favicon.BuildURL("", "https://example.com/path", 64)
	if activeProvider != "self_hosted" {
		t.Fatalf("active provider = %q, want self_hosted", activeProvider)
	}

	var persisted config.FaviconSettings
	if err := dao.GetJsonSetting(db, constant.KeyFaviconProviderConfig, &persisted); err != nil {
		t.Fatalf("read persisted settings: %v", err)
	}
	if persisted.DefaultProviderID != "self_hosted" {
		t.Fatalf("persisted provider = %q", persisted.DefaultProviderID)
	}

	getResponse := performJSONRequest(t, r, http.MethodGet, "/settings", nil)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, body = %s", getResponse.Code, getResponse.Body.String())
	}
	var getBody util.APIResponse[FaviconProviderConfigResponse]
	if err := json.Unmarshal(getResponse.Body.Bytes(), &getBody); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if getBody.Data.DefaultProviderID != "self_hosted" {
		t.Fatalf("GET default provider = %q", getBody.Data.DefaultProviderID)
	}
	if len(getBody.Data.Providers) != 5 {
		t.Fatalf("GET provider count = %d, want 5", len(getBody.Data.Providers))
	}

	invalidResponse := performJSONRequest(t, r, http.MethodPost, "/settings", config.FaviconSettings{
		DefaultProviderID: "missing",
	})
	if invalidResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid status = %d, body = %s", invalidResponse.Code, invalidResponse.Body.String())
	}
	_, activeProvider = favicon.BuildURL("", "https://example.com/path", 64)
	if activeProvider != "self_hosted" {
		t.Fatalf("invalid save changed active provider to %q", activeProvider)
	}

	previewSettings := config.FaviconSettings{
		DefaultProviderID: "preview_only",
		CustomProviders: []config.FaviconProviderConfig{
			{
				ID:          "preview_only",
				Name:        "Preview only",
				URLTemplate: "https://preview.example.test/{host}?size={size}",
				Enabled:     true,
			},
		},
	}
	previewResponse := performJSONRequest(t, r, http.MethodPost, "/preview", FaviconProviderPreviewRequest{
		Settings:   previewSettings,
		ProviderID: "preview_only",
		PageURL:    "https://example.com/path",
		Size:       32,
	})
	if previewResponse.Code != http.StatusOK {
		t.Fatalf("preview status = %d, body = %s", previewResponse.Code, previewResponse.Body.String())
	}
	var previewBody util.APIResponse[FaviconProviderPreviewResponse]
	if err := json.Unmarshal(previewResponse.Body.Bytes(), &previewBody); err != nil {
		t.Fatalf("decode preview response: %v", err)
	}
	if previewBody.Data.URL != "https://preview.example.test/example.com?size=32" {
		t.Fatalf("preview URL = %q", previewBody.Data.URL)
	}
	_, activeProvider = favicon.BuildURL("", "https://example.com/path", 64)
	if activeProvider != "self_hosted" {
		t.Fatalf("preview changed active provider to %q", activeProvider)
	}
}

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()
	var requestBody *bytes.Reader
	if body == nil {
		requestBody = bytes.NewReader(nil)
	} else {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request: %v", err)
		}
		requestBody = bytes.NewReader(jsonBytes)
	}
	req, err := http.NewRequest(method, path, requestBody)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}
