package dao

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SystemSetting struct {
	Key       string `gorm:"primaryKey"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetSetting(db *gorm.DB, key string) (string, error) {
	var s SystemSetting
	err := db.First(&s, "key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return s.Value, nil
}

func SetSetting(db *gorm.DB, key string, value string) error {
	var s SystemSetting
	err := db.First(&s, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s = SystemSetting{
			Key:   key,
			Value: value,
		}
		return db.Create(&s).Error
	}
	if err != nil {
		return err
	}
	s.Value = value
	return db.Save(&s).Error
}

func GetJsonSetting(db *gorm.DB, key string, v interface{}) error {
	val, err := GetSetting(db, key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), v)
}

func SetJsonSetting(db *gorm.DB, key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return SetSetting(db, key, string(bytes))
}
