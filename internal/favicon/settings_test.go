package favicon

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestLoadUsesPersistedSettings(t *testing.T) {
	resetRegistryForTest(t)
	db := newSettingsTestDB(t)
	persisted := config.FaviconSettings{DefaultProviderID: ProviderYandex}
	if err := dao.SetJsonSetting(db, constant.KeyFaviconProviderConfig, persisted); err != nil {
		t.Fatalf("seed settings: %v", err)
	}

	if err := Load(db); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	_, providerID := BuildURL("", "https://example.com", 64)
	if providerID != ProviderYandex {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderYandex)
	}
}

func TestLoadMissingSettingsUsesDefault(t *testing.T) {
	resetRegistryForTest(t)
	db := newSettingsTestDB(t)
	if err := Replace(config.FaviconSettings{DefaultProviderID: ProviderGoogle}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}

	if err := Load(db); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	_, providerID := BuildURL("", "https://example.com", 64)
	if providerID != ProviderGstaticCN {
		t.Fatalf("provider ID = %q, want %q", providerID, ProviderGstaticCN)
	}
}

func TestSavePersistsAndPublishesSnapshot(t *testing.T) {
	resetRegistryForTest(t)
	db := newSettingsTestDB(t)
	settings := config.FaviconSettings{DefaultProviderID: ProviderDuckDuckGo}

	if err := Save(db, settings); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var persisted config.FaviconSettings
	if err := dao.GetJsonSetting(db, constant.KeyFaviconProviderConfig, &persisted); err != nil {
		t.Fatalf("read persisted settings: %v", err)
	}
	if persisted.DefaultProviderID != ProviderDuckDuckGo {
		t.Fatalf("persisted provider = %q", persisted.DefaultProviderID)
	}
	_, providerID := BuildURL("", "https://example.com", 64)
	if providerID != ProviderDuckDuckGo {
		t.Fatalf("active provider = %q", providerID)
	}
}

func TestSaveInvalidSettingsDoesNotPersistOrPublish(t *testing.T) {
	resetRegistryForTest(t)
	db := newSettingsTestDB(t)
	before := Settings()

	err := Save(db, config.FaviconSettings{DefaultProviderID: "missing"})
	if err == nil {
		t.Fatal("Save() error = nil, want validation error")
	}

	var persisted config.FaviconSettings
	if err := dao.GetJsonSetting(db, constant.KeyFaviconProviderConfig, &persisted); err != nil {
		t.Fatalf("read persisted settings: %v", err)
	}
	if persisted.DefaultProviderID != "" {
		t.Fatalf("invalid settings persisted: %+v", persisted)
	}
	if got := Settings(); got.DefaultProviderID != before.DefaultProviderID {
		t.Fatalf("active settings changed: %+v", got)
	}
}

func newSettingsTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&dao.SystemSetting{}); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	return db
}
