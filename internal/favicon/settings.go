package favicon

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/dao"
	"sync"

	"gorm.io/gorm"
)

// Load reads the persisted configuration once and publishes a complete
// immutable snapshot. Missing settings resolve to the built-in defaults.
func Load(db *gorm.DB) error {
	settings := DefaultSettings()
	var persisted config.FaviconSettings
	if err := dao.GetJsonSetting(db, constant.KeyFaviconProviderConfig, &persisted); err != nil {
		return err
	}
	if persisted.DefaultProviderID != "" || len(persisted.CustomProviders) > 0 {
		settings = persisted
	}
	return Replace(settings)
}

// Reload is an explicit hook for future multi-instance invalidation.
func Reload(db *gorm.DB) error {
	return Load(db)
}

// Save validates and compiles before persistence, then atomically publishes
// the exact snapshot that was persisted.
func Save(db *gorm.DB, settings config.FaviconSettings) error {
	compiled, err := compileSnapshot(settings)
	if err != nil {
		return err
	}
	if err := dao.SetJsonSetting(db, constant.KeyFaviconProviderConfig, compiled.settings); err != nil {
		return err
	}
	activeSnapshot.Store(compiled)
	warnedFallbacks = sync.Map{}
	return nil
}
