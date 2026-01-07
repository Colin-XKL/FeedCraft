package dao

import "gorm.io/gorm"

type Blueprint struct {
	BaseModelWithoutPK
	Name            string          `gorm:"primaryKey" json:"name" binding:"required"`
	Description     string          `json:"description,omitempty"`
	BlueprintConfig []BlueprintItem `json:"blueprint_config" binding:"required" gorm:"serializer:json"`
}

func (Blueprint) TableName() string {
	return "blueprints"
}

type BlueprintItem struct {
	ProcessorName string            `json:"processor_name" binding:"required"` // Formerly CraftName
	Option        map[string]string `json:"option"`
}

// CreateBlueprint creates a new Blueprint record
func CreateBlueprint(db *gorm.DB, blueprint *Blueprint) error {
	return db.Create(blueprint).Error
}

// GetBlueprintByName retrieves a Blueprint record by its Name
func GetBlueprintByName(db *gorm.DB, name string) (*Blueprint, error) {
	var blueprint Blueprint
	result := db.Where("name = ?", name).First(&blueprint)
	return &blueprint, result.Error
}

// GetAllBlueprints retrieves all Blueprint records
func GetAllBlueprints(db *gorm.DB) ([]Blueprint, error) {
	var blueprints []Blueprint
	result := db.Find(&blueprints)
	return blueprints, result.Error
}

// UpdateBlueprint updates an existing Blueprint record
func UpdateBlueprint(db *gorm.DB, blueprint *Blueprint) error {
	return db.Save(blueprint).Error
}

// DeleteBlueprint deletes a Blueprint record by its Name
func DeleteBlueprint(db *gorm.DB, name string) error {
	blueprint, err := GetBlueprintByName(db, name)
	if err != nil {
		return err
	}
	return db.Delete(blueprint).Error
}

// ListBlueprints retrieves all Blueprint records
func ListBlueprints(db *gorm.DB) ([]Blueprint, error) {
	var blueprints []Blueprint
	result := db.Find(&blueprints)
	if result.Error != nil {
		return nil, result.Error
	}
	return blueprints, nil
}