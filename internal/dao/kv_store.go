package dao

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type KVStore struct {
	Key       string `gorm:"primaryKey"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetValue(db *gorm.DB, key string) (string, error) {
	var kv KVStore
	err := db.First(&kv, "key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return kv.Value, nil
}

func SetValue(db *gorm.DB, key string, value string) error {
	var kv KVStore
	err := db.First(&kv, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		kv = KVStore{
			Key:   key,
			Value: value,
		}
		return db.Create(&kv).Error
	}
	if err != nil {
		return err
	}
	kv.Value = value
	return db.Save(&kv).Error
}

// GetJsonValue unmarshals a JSON value into v
func GetJsonValue(db *gorm.DB, key string, v interface{}) error {
	val, err := GetValue(db, key)
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	return json.Unmarshal([]byte(val), v)
}

// SetJsonValue marshals v to JSON and saves it
func SetJsonValue(db *gorm.DB, key string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return SetValue(db, key, string(bytes))
}
