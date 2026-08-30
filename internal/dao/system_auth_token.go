package dao

import (
	"time"

	"gorm.io/gorm"
)

type SystemAuthToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	Label     string    `json:"label,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (SystemAuthToken) TableName() string {
	return "system_auth_tokens"
}

func CreateSystemAuthToken(db *gorm.DB, token *SystemAuthToken) error {
	return db.Create(token).Error
}

func GetSystemAuthTokenByID(db *gorm.DB, id uint) (*SystemAuthToken, error) {
	var token SystemAuthToken
	result := db.Where("id = ?", id).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func GetSystemAuthTokenByToken(db *gorm.DB, value string) (*SystemAuthToken, error) {
	var token SystemAuthToken
	if value == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("token = ?", value).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func DeleteSystemAuthToken(db *gorm.DB, id uint) error {
	var token SystemAuthToken
	result := db.Where("id = ?", id).Delete(&token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func ListSystemAuthTokens(db *gorm.DB) ([]*SystemAuthToken, error) {
	var tokens []*SystemAuthToken
	if err := db.Order("created_at desc").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}
