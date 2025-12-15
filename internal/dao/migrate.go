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
	err := db.AutoMigrate(
		&CustomRecipe{},
		&CustomRecipeV2{}, // Create the new V2 table
		&CraftFlow{}, &CraftAtom{},
		&User{}, // 确保 User 表被初始化
	)
	if err != nil {
		logrus.Error("migrate database error.", err)
		return
	}

	// Perform data migration from custom_recipes to custom_recipes_v2
	migrateRecipesToV2(db)

	logrus.Info("migrate database done.")

	// 创建 admin 账户
	createAdminUser(db)
}

func migrateRecipesToV2(db *gorm.DB) {
	if !db.Migrator().HasTable(&CustomRecipe{}) {
		logrus.Info("original recipe table does not exist, skipping migration.")
		return
	}

	logrus.Info("starting migration from 'custom_recipes' to 'custom_recipes_v2'...")

	var oldRecipes []*CustomRecipe
	if err := db.Find(&oldRecipes).Error; err != nil {
		logrus.Errorf("failed to query old recipes for migration: %v", err)
		return
	}

	for _, oldR := range oldRecipes {
		// Check if a recipe with the same ID already exists in the V2 table.
		var existingV2 CustomRecipeV2
		if err := db.First(&existingV2, "id = ?", oldR.ID).Error; err == nil {
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

		newRecipeV2 := CustomRecipeV2{
			ID:           oldR.ID,
			Description:  oldR.Description,
			Craft:        oldR.Craft,
			SourceType:   string(constant.SourceRSS), // Store as string in DB
			SourceConfig: string(configJSON),
		}

		if err := db.Create(&newRecipeV2).Error; err != nil {
			logrus.Errorf("failed to insert V2 recipe for id %s: %v", oldR.ID, err)
		}
	}

	logrus.Info("recipe migration to v2 completed.")
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
