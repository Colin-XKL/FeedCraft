package dao

import (
	"FeedCraft/internal/config"
	"FeedCraft/internal/constant"
	"FeedCraft/internal/util"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrateDatabases() {
	logrus.Info("migrating databases...")
	db := util.GetDatabase()

	// 1. Rename tables if they exist with old names
	renameTables(db)

	// 1.1 Rename columns if necessary
	renameColumns(db)

	// 2. AutoMigrate with new structs
	err := db.AutoMigrate(
		&CustomRecipe{}, // Keep for legacy data migration source
		&Channel{},      // New Channel table
		&Blueprint{},    // New Blueprint table
		&Tool{},         // New Tool table
		&User{},
		&SystemSetting{},
	)
	if err != nil {
		logrus.Error("migrate database error.", err)
		return
	}

	// 3. Perform data migration from legacy custom_recipes to channels
	migrateRecipesToChannels(db)

	logrus.Info("migrate database done.")

	// Create admin account
	createAdminUser(db)
}

func renameTables(db *gorm.DB) {
	// Rename craft_atoms -> tools
	if db.Migrator().HasTable("craft_atoms") && !db.Migrator().HasTable("tools") {
		logrus.Info("Renaming table 'craft_atoms' to 'tools'")
		if err := db.Migrator().RenameTable("craft_atoms", "tools"); err != nil {
			logrus.Errorf("Failed to rename table 'craft_atoms': %v", err)
		}
	}

	// Rename craft_flows -> blueprints
	if db.Migrator().HasTable("craft_flows") && !db.Migrator().HasTable("blueprints") {
		logrus.Info("Renaming table 'craft_flows' to 'blueprints'")
		if err := db.Migrator().RenameTable("craft_flows", "blueprints"); err != nil {
			logrus.Errorf("Failed to rename table 'craft_flows': %v", err)
		}
	}

	// Rename custom_recipes_v2 -> channels
	if db.Migrator().HasTable("custom_recipes_v2") && !db.Migrator().HasTable("channels") {
		logrus.Info("Renaming table 'custom_recipes_v2' to 'channels'")
		if err := db.Migrator().RenameTable("custom_recipes_v2", "channels"); err != nil {
			logrus.Errorf("Failed to rename table 'custom_recipes_v2': %v", err)
		}
	}
}

func renameColumns(db *gorm.DB) {
	// Rename channels.craft -> channels.processor_name
	if db.Migrator().HasTable("channels") && db.Migrator().HasColumn("channels", "craft") {
		logrus.Info("Renaming column 'channels.craft' to 'processor_name'")
		if err := db.Migrator().RenameColumn("channels", "craft", "processor_name"); err != nil {
			logrus.Errorf("Failed to rename column 'channels.craft': %v", err)
		}
	}

	// Rename blueprints.flow_config -> blueprints.blueprint_config
	if db.Migrator().HasTable("blueprints") && db.Migrator().HasColumn("blueprints", "flow_config") {
		logrus.Info("Renaming column 'blueprints.flow_config' to 'blueprint_config'")
		if err := db.Migrator().RenameColumn("blueprints", "flow_config", "blueprint_config"); err != nil {
			logrus.Errorf("Failed to rename column 'blueprints.flow_config': %v", err)
		}
	}
}

func migrateRecipesToChannels(db *gorm.DB) {
	if !db.Migrator().HasTable(&CustomRecipe{}) {
		logrus.Info("original recipe table does not exist, skipping migration.")
		return
	}

	logrus.Info("starting migration from 'custom_recipes' to 'channels'...")

	var oldRecipes []*CustomRecipe
	if err := db.Find(&oldRecipes).Error; err != nil {
		logrus.Errorf("failed to query old recipes for migration: %v", err)
		return
	}

	for _, oldR := range oldRecipes {
		// Check if a channel with the same ID already exists.
		var existingChannel Channel
		if err := db.First(&existingChannel, "id = ?", oldR.ID).Error; err == nil {
			// Record already exists, skip.
			continue
		}

		if oldR.FeedURL == "" {
			continue
		}

		// Create the new nested SourceConfig structure
		newConfig := config.SourceConfig{
			Type: constant.SourceRSS,
			HttpFetcher: &config.HttpFetcherConfig{
				URL: oldR.FeedURL,
			},
		}

		configJSON, err := json.Marshal(newConfig)
		if err != nil {
			logrus.Errorf("failed to marshal new source config for recipe id %s: %v", oldR.ID, err)
			continue
		}

		newChannel := Channel{
			ID:            oldR.ID,
			Description:   oldR.Description,
			ProcessorName: oldR.Craft,                  // Map Craft -> ProcessorName
			SourceType:    string(constant.SourceRSS), // Store as string in DB
			SourceConfig:  string(configJSON),
		}

		if err := db.Create(&newChannel).Error; err != nil {
			logrus.Errorf("failed to insert Channel for id %s: %v", oldR.ID, err)
		}
	}

	logrus.Info("recipe migration to channels completed.")
}

var defaultAdminUsername = "admin"
var defaultPassword = "adminadmin" // default defaultPassword string

var defaultAdminUser = User{
	Username: defaultAdminUsername,
	NickName: "Admin",
	Email:    "admin@example.com",
}

func createAdminUser(db *gorm.DB) {
	md5Password := util.GetMD5Hash(defaultPassword)

	// 检查是否已经存在 admin 用户
	var user User
	result := db.Where("username = ?", defaultAdminUsername).First(&user)
	if result.Error == nil {
		logrus.Info("admin user already exists")
		return
	}

	// 创建 admin 用户
	if err := CreateUser(db, &defaultAdminUser, md5Password); err != nil {
		logrus.Error("failed to create admin user:", err)
		return
	}

	logrus.Info("admin user created successfully")
}

// 重置 admin 密码
func ResetAdminPassword() error {
	logrus.Info("resetting admin password...")
	db := util.GetDatabase()
	md5Password := util.GetMD5Hash(defaultPassword)
	return UpdateUserPassword(db, &defaultAdminUser, md5Password)
}