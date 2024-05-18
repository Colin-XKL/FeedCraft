package dao

import (
	"FeedCraft/internal/util"
	"github.com/sirupsen/logrus"
)

func MigrateDatabases() {
	logrus.Info("migrating databases...")
	db := util.GetDatabase()
	err := db.AutoMigrate(
		&CustomRecipe{},
	)
	if err != nil {
		logrus.Error("migrate database error.", err)
		return
	}
	logrus.Info("migrate database done.")
}
