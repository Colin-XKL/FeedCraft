package util

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const sqliteDbName = "feed-craft.db"

var (
	db   *gorm.DB
	once sync.Once
)

func GetDatabase() *gorm.DB {
	once.Do(func() {
		envClient := GetEnvClient()
		sqlitePath := envClient.GetString("DB_SQLITE_PATH") // 这里path 实际为文件夹,后续需要在文档中说明清楚
		if sqlitePath == "" {
			log.Fatalf("DB_SQLITE_PATH not valid")
		}
		path := filepath.Join(sqlitePath, sqliteDbName)
		conf := &gorm.Config{}
		// use a pure go version sqlite driver to make it more portable
		var err error
		db, err = gorm.Open(sqlite.Open(path), conf)
		if err != nil || db == nil {
			logrus.Fatalf("failed to connect to database. %s", err)
		}
		// 设置连接池
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Errorf("failed to get underlying sql.DB: %v", err)
		} else {
			// SetMaxIdleConns 设置空闲连接池中连接的最大数量
			sqlDB.SetMaxIdleConns(10)
			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(100)
			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(time.Hour)
		}
	})
	return db
}

func CloseDatabase() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
