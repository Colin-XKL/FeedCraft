package dao

import (
	"FeedCraft/internal/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	password := "adminadmin"

	// 检查是否已经存在 admin 用户
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error == nil {
		logrus.Info("admin user already exists")
		return
	}

	// 创建 admin 用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("failed to hash password for admin user", err)
		return
	}

	adminUser := User{
		Username:     username,
		NickName:     "Admin",
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
	}

	if err := db.Create(&adminUser).Error; err != nil {
		logrus.Error("failed to create admin user", err)
		return
	}

	logrus.Info("admin user created successfully")
}
