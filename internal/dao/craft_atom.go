package dao

// CraftAtom craft atom 由 craft template 派生出来
type CraftAtom struct {
	BaseModelWithPK
	Name         string            `json:"name" gorm:"primaryKey" binding:"required"`
	Description  string            `json:"description"`
	TemplateName string            `json:"template_name" binding:"required" `
	Params       map[string]string `json:"params" gorm:"serializer:json" gorm:"index"`
}
