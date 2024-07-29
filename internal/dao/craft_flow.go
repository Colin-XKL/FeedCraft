package dao

import "gorm.io/gorm"

type CraftFlow struct {
	BaseModelWithoutPK
	Name            string          `gorm:"primaryKey" json:"name" binding:"required"`
	Description     string          `json:"description,omitempty"`
	CraftFlowConfig []CraftFlowItem `json:"craft_flow_config" binding:"required" gorm:"serializer:json"`
}

type CraftFlowItem struct {
	CraftName string            `json:"craft_name" binding:"required"`
	Option    map[string]string `json:"option"`
}

// CreateCraftFlow creates a new CraftFlow record
func CreateCraftFlow(db *gorm.DB, craftFlow *CraftFlow) error {
	return db.Create(craftFlow).Error
}

// GetCraftFlowByName retrieves a CraftFlow record by its Name
func GetCraftFlowByName(db *gorm.DB, name string) (*CraftFlow, error) {
	var craftFlow CraftFlow
	result := db.Where("name = ?", name).First(&craftFlow)
	return &craftFlow, result.Error
}

// GetAllCraftFlows retrieves all CraftFlow records
func GetAllCraftFlows(db *gorm.DB) ([]CraftFlow, error) {
	var craftFlows []CraftFlow
	result := db.Find(&craftFlows)
	return craftFlows, result.Error
}

// UpdateCraftFlow updates an existing CraftFlow record
func UpdateCraftFlow(db *gorm.DB, craftFlow *CraftFlow) error {
	return db.Save(craftFlow).Error
}

// DeleteCraftFlow deletes a CraftFlow record by its Name
func DeleteCraftFlow(db *gorm.DB, name string) error {
	craftFlow, err := GetCraftFlowByName(db, name)
	if err != nil {
		return err
	}
	return db.Delete(craftFlow).Error
}

// ListCraftFlows retrieves all CraftFlow records
func ListCraftFlows(db *gorm.DB) ([]CraftFlow, error) {
	var craftFlows []CraftFlow
	result := db.Find(&craftFlows)
	if result.Error != nil {
		return nil, result.Error
	}
	return craftFlows, nil
}
