package dao

import (
	"FeedCraft/internal/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrateDatabases() {
	logrus.Info("migrating databases...")
	db := util.GetDatabase()
	err := db.AutoMigrate(
		&CustomRecipe{},
		&CraftFlow{}, &CraftAtom{},
		&User{}, // 确保 User 表被初始化
	)
	if err != nil {
		logrus.Error("migrate database error.", err)
		return
	}
	logrus.Info("migrate database done.")

	// 创建 admin 账户
	createAdminUser(db)
}

func createAdminUser(db *gorm.DB) {
	username := "admin"
	defaultPassword := "adminadmin" // default defaultPassword string
	md5Password := util.GetMD5Hash(defaultPassword)

	// 检查是否已经存在 admin 用户
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error == nil {
		logrus.Info("admin user already exists")
		return
	}

	// 创建 admin 用户
	adminUser := &User{
		Username: username,
		NickName: "Admin",
		Email:    "admin@example.com",
	}
	if err := CreateUser(db, adminUser, md5Password); err != nil {
		logrus.Error("failed to create admin user:", err)
		return
	}

	logrus.Info("admin user created successfully")
}
