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
		// &CustomRecipe{},
		&CustomRecipeV2{}, // Create the new V2 table
		&CraftFlow{}, &CraftAtom{},
		&TopicFeed{}, // Add TopicFeed migration
		&User{},      // 确保 User 表被初始化
		&SystemSetting{},
		&ExecutionLog{},
		&ResourceHealth{},
		&SystemNotification{},
		&Inbox{},
		&InboxItem{},
		&SystemAuthToken{},
	)
	if err != nil {
		logrus.Fatalf("migrate database error: %v", err)
		return
	}

	// Perform data migration from custom_recipes to custom_recipes_v2
	migrateRecipesToV2(db)
	migrateTopicFeedInputs(db)

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

func migrateTopicFeedInputs(db *gorm.DB) {
	if !db.Migrator().HasTable(&TopicFeed{}) {
		return
	}
	if !db.Migrator().HasColumn("topic_feeds", "input_uris") {
		return
	}

	type legacyTopicInputRow struct {
		ID        string
		Inputs    string
		InputURIs string
	}

	var rows []legacyTopicInputRow
	if err := db.Table("topic_feeds").Select("id, inputs, input_uris").Scan(&rows).Error; err != nil {
		logrus.Errorf("failed to scan topic feed input_uris migration rows: %v", err)
		return
	}

	for _, row := range rows {
		if row.InputURIs == "" || row.InputURIs == "null" || row.InputURIs == "[]" {
			continue
		}
		if row.Inputs != "" && row.Inputs != "null" && row.Inputs != "[]" {
			continue
		}

		var uris []string
		if err := json.Unmarshal([]byte(row.InputURIs), &uris); err != nil {
			logrus.Errorf("failed to unmarshal legacy input_uris for topic %s: %v", row.ID, err)
			continue
		}

		inputs := make([]TopicInput, 0, len(uris))
		for _, uri := range uris {
			inputs = append(inputs, TopicInput{URI: uri})
		}
		inputsJSON, err := json.Marshal(inputs)
		if err != nil {
			logrus.Errorf("failed to marshal migrated inputs for topic %s: %v", row.ID, err)
			continue
		}
		if err := db.Table("topic_feeds").Where("id = ?", row.ID).Update("inputs", string(inputsJSON)).Error; err != nil {
			logrus.Errorf("failed to migrate topic inputs for topic %s: %v", row.ID, err)
		}
	}

	if err := db.Exec("ALTER TABLE topic_feeds DROP COLUMN input_uris").Error; err != nil {
		logrus.Errorf("failed to drop legacy topic input_uris column: %v", err)
		return
	}
	logrus.Info("topic feed input_uris migration completed.")
}

var defaultAdminUsername = "admin"
var defaultPassword = "adminadmin" // default defaultPassword string

var defaultAdminUser = User{
	Username: defaultAdminUsername,
	NickName: "Admin",
	Email:    "admin@example.com",
}

func createAdminUser(db *gorm.DB) {
	md5Password := util.GetPasswordMD5Hash(defaultPassword)

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
	md5Password := util.GetPasswordMD5Hash(defaultPassword)
	return UpdateUserPassword(db, &defaultAdminUser, md5Password)
}
