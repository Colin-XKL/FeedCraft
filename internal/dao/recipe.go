package dao

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type BaseModelWithoutPK struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`
}

// CustomRecipe represents the original recipe structure.
type CustomRecipe struct {
	BaseModelWithoutPK
	ID          string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Craft       string `json:"craft" binding:"required"`
	FeedURL     string `json:"feed_url" binding:"required"`
}

// CustomRecipeV2 represents the new, refactored recipe structure.
type CustomRecipeV2 struct {
	BaseModelWithoutPK
	ID           string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
	Description  string `json:"description,omitempty"`
	Craft        string `json:"craft" binding:"required"`
	SourceType   string `gorm:"type:varchar(50);not null;default:'rss'" json:"source_type"`
	SourceConfig string `gorm:"type:text" json:"source_config"`
}

func (CustomRecipeV2) TableName() string {
	return "custom_recipes_v2"
}

// CreateCustomRecipe creates a new CustomRecipe record
func CreateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Create(recipe).Error
}

// CreateCustomRecipeV2 creates a new CustomRecipeV2 record
func CreateCustomRecipeV2(db *gorm.DB, recipe *CustomRecipeV2) error {
	return db.Create(recipe).Error
}

// GetCustomRecipeByID retrieves a CustomRecipe record by its ID
func GetCustomRecipeByID(db *gorm.DB, id string) (*CustomRecipe, error) {
	var recipe CustomRecipe
	// 添加日志记录查询的 ID
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&recipe)
	return &recipe, result.Error
}

// GetCustomRecipeByIDV2 retrieves a V2 CustomRecipe record by its ID
func GetCustomRecipeByIDV2(db *gorm.DB, id string) (*CustomRecipeV2, error) {
	var recipe CustomRecipeV2
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&recipe)
	return &recipe, result.Error
}

// UpdateCustomRecipe updates an existing CustomRecipe record
func UpdateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Save(recipe).Error
}

// UpdateCustomRecipeV2 updates an existing CustomRecipeV2 record
func UpdateCustomRecipeV2(db *gorm.DB, recipe *CustomRecipeV2) error {
	return db.Save(recipe).Error
}

// DeleteCustomRecipe deletes a CustomRecipe record by its ID
func DeleteCustomRecipe(db *gorm.DB, id string) error {
	var recipe CustomRecipe
	return db.Where("id = ?", id).Delete(&recipe).Error
}

// DeleteCustomRecipeV2 deletes a CustomRecipeV2 record by its ID
func DeleteCustomRecipeV2(db *gorm.DB, id string) error {
	var recipe CustomRecipeV2
	return db.Where("id = ?", id).Delete(&recipe).Error
}

func ListCustomRecipe(db *gorm.DB) ([]*CustomRecipe, error) {
	var recipes []*CustomRecipe
	if err := db.Find(&recipes).Error; err != nil {
		return recipes, err
	}
	return recipes, nil
}

func ListCustomRecipeV2(db *gorm.DB) ([]*CustomRecipeV2, error) {
	var recipes []*CustomRecipeV2
	if err := db.Find(&recipes).Error; err != nil {
		return recipes, err
	}
	return recipes, nil
}
