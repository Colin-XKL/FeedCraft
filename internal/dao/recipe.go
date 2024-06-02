package dao

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type BaseModelWithPK struct {
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`
}

type CustomRecipe struct {
	BaseModelWithPK
	ID          string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Craft       string `json:"craft" binding:"required"`
	FeedURL     string `json:"feed_url" binding:"required"`
}

// CreateCustomRecipe creates a new CustomRecipe record
func CreateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Create(recipe).Error
}

// GetCustomRecipeByID retrieves a CustomRecipe record by its ID
func GetCustomRecipeByID(db *gorm.DB, id string) (*CustomRecipe, error) {
	var recipe CustomRecipe
	result := db.Where("id = ?", id).First(&recipe)
	return &recipe, result.Error
}

// UpdateCustomRecipe updates an existing CustomRecipe record
func UpdateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Save(recipe).Error
}

// DeleteCustomRecipe deletes a CustomRecipe record by its ID
func DeleteCustomRecipe(db *gorm.DB, id string) error {
	var recipe CustomRecipe
	return db.Where("id = ?", id).Delete(&recipe).Error
}

func ListCustomRecipe(db *gorm.DB) ([]*CustomRecipe, error) {
	var recipes []*CustomRecipe
	if err := db.Find(&recipes).Error; err != nil {
		return recipes, err
	}
	return recipes, nil
}
