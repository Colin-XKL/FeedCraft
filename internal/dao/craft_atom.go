package dao

import "gorm.io/gorm"

// CraftAtom craft atom 由 craft template 派生出来
type CraftAtom struct {
	BaseModelWithoutPK
	Name         string            `json:"name" gorm:"primaryKey" binding:"required"`
	Description  string            `json:"description"`
	TemplateName string            `json:"template_name" binding:"required" `
	Params       map[string]string `json:"params" gorm:"serializer:json;index"`
}

// CreateCraftAtom creates a new CraftAtom record
func CreateCraftAtom(db *gorm.DB, craftAtom *CraftAtom) error {
	return db.Create(craftAtom).Error
}

// GetCraftAtomByName retrieves a CraftAtom record by its Name
func GetCraftAtomByName(db *gorm.DB, name string) (*CraftAtom, error) {
	var craftAtom CraftAtom
	result := db.Where("name = ?", name).First(&craftAtom)
	return &craftAtom, result.Error
}

// GetAllCraftAtoms retrieves all CraftAtom records
func GetAllCraftAtoms(db *gorm.DB) ([]CraftAtom, error) {
	var craftAtoms []CraftAtom
	result := db.Find(&craftAtoms)
	return craftAtoms, result.Error
}

// UpdateCraftAtom updates an existing CraftAtom record
func UpdateCraftAtom(db *gorm.DB, craftAtom *CraftAtom) error {
	return db.Save(craftAtom).Error
}

// DeleteCraftAtom deletes a CraftAtom record by its Name
func DeleteCraftAtom(db *gorm.DB, name string) error {
	craftAtom, err := GetCraftAtomByName(db, name)
	if err != nil {
		return err
	}
	return db.Delete(craftAtom).Error
}
