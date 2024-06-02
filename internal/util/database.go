package util

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"path/filepath"
)

const sqliteDbName = "feed-craft.db"

func GetDatabase() *gorm.DB {
	envClient := GetEnvClient()
	sqlitePath := envClient.GetString("DB_SQLITE_PATH")
	if sqlitePath == "" {
		log.Fatalf("DB_SQLITE_PATH not valid")
	}

	path := filepath.Join(sqlitePath, sqliteDbName)

	conf := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open(path), conf)
	if err != nil || db == nil {
		logrus.Fatalf("failed to connect to database. %s", err)
		panic("failed to connect database")
	}
	return db
}
