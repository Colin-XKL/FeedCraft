package dao

import (
	"time"

	"gorm.io/gorm"
)

type Inbox struct {
	BaseModelWithoutPK
	ID          string `gorm:"primaryKey" json:"id" binding:"required"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MaxItems    int    `gorm:"default:100" json:"max_items"`
	IsPublic    bool   `gorm:"default:true" json:"is_public"`
}

type InboxItem struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	InboxID     string    `gorm:"uniqueIndex:idx_inbox_item_id;not null" json:"inbox_id"`
	ItemID      string    `gorm:"uniqueIndex:idx_inbox_item_id;not null" json:"item_id"`
	Title       string    `gorm:"not null" json:"title"`
	URL         string    `json:"url,omitempty"`
	Content     string    `gorm:"type:text" json:"content,omitempty"`
	Summary     string    `json:"summary,omitempty"`
	Author      string    `json:"author,omitempty"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Inbox) TableName() string {
	return "inboxes"
}

func (InboxItem) TableName() string {
	return "inbox_items"
}

func CreateInbox(db *gorm.DB, inbox *Inbox) error {
	return db.Create(inbox).Error
}

func GetInboxByID(db *gorm.DB, id string) (*Inbox, error) {
	var inbox Inbox
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&inbox)
	if result.Error != nil {
		return nil, result.Error
	}
	return &inbox, nil
}

func UpdateInbox(db *gorm.DB, inbox *Inbox) error {
	return db.Save(inbox).Error
}

func DeleteInbox(db *gorm.DB, id string) error {
	var inbox Inbox
	result := db.Where("id = ?", id).Delete(&inbox)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func ListInboxes(db *gorm.DB) ([]*Inbox, error) {
	var inboxes []*Inbox
	if err := db.Order("created_at desc").Find(&inboxes).Error; err != nil {
		return nil, err
	}
	return inboxes, nil
}

func GetInboxItemByItemID(db *gorm.DB, inboxID, itemID string) (*InboxItem, error) {
	var item InboxItem
	if inboxID == "" || itemID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("inbox_id = ? AND item_id = ?", inboxID, itemID).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func ListInboxItems(db *gorm.DB, inboxID string) ([]*InboxItem, error) {
	var items []*InboxItem
	if inboxID == "" {
		return items, nil
	}
	if err := db.Where("inbox_id = ?", inboxID).Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func CountInboxItems(db *gorm.DB, inboxID string) (int64, error) {
	var count int64
	if err := db.Model(&InboxItem{}).Where("inbox_id = ?", inboxID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func DeleteOverflowInboxItems(db *gorm.DB, inboxID string, maxItems int) error {
	if maxItems <= 0 {
		return db.Where("inbox_id = ?", inboxID).Delete(&InboxItem{}).Error
	}

	var ids []uint
	if err := db.Model(&InboxItem{}).
		Where("inbox_id = ?", inboxID).
		Order("created_at desc").
		Offset(maxItems).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}
	return db.Where("id IN ?", ids).Delete(&InboxItem{}).Error
}

func ListExistingInboxItemIDs(db *gorm.DB, inboxID string, itemIDs []string) ([]string, error) {
	if len(itemIDs) == 0 {
		return nil, nil
	}
	var existing []string
	if err := db.Model(&InboxItem{}).
		Where("inbox_id = ? AND item_id IN ?", inboxID, itemIDs).
		Pluck("item_id", &existing).Error; err != nil {
		return nil, err
	}
	return existing, nil
}
