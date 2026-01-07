package dao

import (
	"gorm.io/gorm"
)

// CustomRecipe represents the original recipe structure.
// Deprecated: Use Channel instead. Kept for migration.
type CustomRecipe struct {
	BaseModelWithoutPK
	ID          string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
	Description string `json:"description,omitempty"`
	Craft       string `json:"craft" binding:"required"`
	FeedURL     string `json:"feed_url" binding:"required"`
}

func (CustomRecipe) TableName() string {
	return "custom_recipes"
}

// Channel (formerly CustomRecipeV2) represents the persisted subscription configuration.
type Channel struct {
	BaseModelWithoutPK
	ID            string `gorm:"primaryKey" json:"id,omitempty" binding:"required"`
	Description   string `json:"description,omitempty"`
	ProcessorName string `json:"processor_name" binding:"required"` // Formerly Craft
	SourceType    string `gorm:"type:varchar(50);not null;default:'rss'" json:"source_type"`
	SourceConfig  string `gorm:"type:text" json:"source_config"`
}

func (Channel) TableName() string {
	return "channels"
}

// CreateCustomRecipe creates a new CustomRecipe record
// Deprecated
func CreateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Create(recipe).Error
}

// CreateChannel creates a new Channel record
func CreateChannel(db *gorm.DB, channel *Channel) error {
	return db.Create(channel).Error
}

// GetCustomRecipeByID retrieves a CustomRecipe record by its ID
// Deprecated
func GetCustomRecipeByID(db *gorm.DB, id string) (*CustomRecipe, error) {
	var recipe CustomRecipe
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&recipe)
	return &recipe, result.Error
}

// GetChannelByID retrieves a Channel record by its ID
func GetChannelByID(db *gorm.DB, id string) (*Channel, error) {
	var channel Channel
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&channel)
	return &channel, result.Error
}

// UpdateCustomRecipe updates an existing CustomRecipe record
// Deprecated
func UpdateCustomRecipe(db *gorm.DB, recipe *CustomRecipe) error {
	return db.Save(recipe).Error
}

// UpdateChannel updates an existing Channel record
func UpdateChannel(db *gorm.DB, channel *Channel) error {
	return db.Save(channel).Error
}

// DeleteCustomRecipe deletes a CustomRecipe record by its ID
// Deprecated
func DeleteCustomRecipe(db *gorm.DB, id string) error {
	var recipe CustomRecipe
	return db.Where("id = ?", id).Delete(&recipe).Error
}

// DeleteChannel deletes a Channel record by its ID
func DeleteChannel(db *gorm.DB, id string) error {
	var channel Channel
	return db.Where("id = ?", id).Delete(&channel).Error
}

// Deprecated
func ListCustomRecipe(db *gorm.DB) ([]*CustomRecipe, error) {
	var recipes []*CustomRecipe
	if err := db.Find(&recipes).Error; err != nil {
		return recipes, err
	}
	return recipes, nil
}

func ListChannels(db *gorm.DB) ([]*Channel, error) {
	var channels []*Channel
	if err := db.Find(&channels).Error; err != nil {
		return channels, err
	}
	return channels, nil
}