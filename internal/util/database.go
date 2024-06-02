package util

import (
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
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
	// use a pure go version sqlite driver to make it more portable
	db, err := gorm.Open(sqlite.Open(path), conf)
	if err != nil || db == nil {
		logrus.Fatalf("failed to connect to database. %s", err)
		panic("failed to connect database")
	}
	return db
}
