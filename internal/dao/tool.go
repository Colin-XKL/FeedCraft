package dao

import "gorm.io/gorm"

// Tool (formerly CraftAtom) derived from craft template
type Tool struct {
	BaseModelWithoutPK
	Name         string            `json:"name" gorm:"primaryKey" binding:"required"`
	Description  string            `json:"description"`
	TemplateName string            `json:"template_name" binding:"required" `
	Params       map[string]string `json:"params" gorm:"serializer:json;index"`
}

func (Tool) TableName() string {
	return "tools"
}

// CreateTool creates a new Tool record
func CreateTool(db *gorm.DB, tool *Tool) error {
	return db.Create(tool).Error
}

// GetToolByName retrieves a Tool record by its Name
func GetToolByName(db *gorm.DB, name string) (*Tool, error) {
	var tool Tool
	result := db.Where("name = ?", name).First(&tool)
	return &tool, result.Error
}

// GetAllTools retrieves all Tool records
func GetAllTools(db *gorm.DB) ([]Tool, error) {
	var tools []Tool
	result := db.Find(&tools)
	return tools, result.Error
}

// UpdateTool updates an existing Tool record
func UpdateTool(db *gorm.DB, tool *Tool) error {
	return db.Save(tool).Error
}

// DeleteTool deletes a Tool record by its Name
func DeleteTool(db *gorm.DB, name string) error {
	tool, err := GetToolByName(db, name)
	if err != nil {
		return err
	}
	return db.Delete(tool).Error
}