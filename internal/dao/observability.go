package dao

import (
	"time"

	"gorm.io/gorm"
)

const (
	ResourceTypeRecipe = "recipe"
	ResourceTypeTopic  = "topic"
)

const (
	ResourceStatusHealthy  = "healthy"
	ResourceStatusDegraded = "degraded"
	ResourceStatusPaused   = "paused"
)

const (
	ExecutionStatusSuccess        = "success"
	ExecutionStatusPartialSuccess = "partial_success"
	ExecutionStatusFailure        = "failure"
	ExecutionStatusPausedSkip     = "paused_skip"
)

type ExecutionLog struct {
	ID uint `gorm:"primaryKey" json:"id"`
	BaseModelWithoutPK
	ResourceType string `gorm:"type:varchar(32);index:idx_exec_resource_time,priority:1;not null" json:"resource_type"`
	ResourceID   string `gorm:"type:varchar(255);index:idx_exec_resource_time,priority:2;not null" json:"resource_id"`
	ResourceName string `gorm:"type:varchar(255)" json:"resource_name"`
	Trigger      string `gorm:"type:varchar(32);index;not null" json:"trigger"`
	Status       string `gorm:"type:varchar(32);index;not null" json:"status"`
	ErrorKind    string `gorm:"type:varchar(64);index" json:"error_kind"`
	Message      string `gorm:"type:text" json:"message"`
	DetailsJSON  string `gorm:"type:text" json:"details_json"`
	RequestID    string `gorm:"type:varchar(128);index" json:"request_id"`
	DurationMS   int64  `json:"duration_ms"`
}

func (ExecutionLog) TableName() string {
	return "execution_logs"
}

type ResourceHealth struct {
	BaseModelWithoutPK
	ResourceType        string     `gorm:"type:varchar(32);primaryKey" json:"resource_type"`
	ResourceID          string     `gorm:"type:varchar(255);primaryKey" json:"resource_id"`
	ResourceName        string     `gorm:"type:varchar(255)" json:"resource_name"`
	CurrentStatus       string     `gorm:"type:varchar(32);index;not null" json:"current_status"`
	ConsecutiveFailures int        `gorm:"not null;default:0" json:"consecutive_failures"`
	LastSuccessAt       *time.Time `json:"last_success_at"`
	LastFailureAt       *time.Time `json:"last_failure_at"`
	LastErrorKind       string     `gorm:"type:varchar(64)" json:"last_error_kind"`
	LastErrorMessage    string     `gorm:"type:text" json:"last_error_message"`
	PausedAt            *time.Time `json:"paused_at"`
	PausedReason        string     `gorm:"type:text" json:"paused_reason"`
}

func (ResourceHealth) TableName() string {
	return "resource_health"
}

type SystemNotification struct {
	ID uint `gorm:"primaryKey" json:"id"`
	BaseModelWithoutPK
	ResourceType string `gorm:"type:varchar(32);index;not null" json:"resource_type"`
	ResourceID   string `gorm:"type:varchar(255);index;not null" json:"resource_id"`
	EventType    string `gorm:"type:varchar(64);index;not null" json:"event_type"`
	Title        string `gorm:"type:varchar(255);not null" json:"title"`
	Content      string `gorm:"type:text" json:"content"`
	StatusAfter  string `gorm:"type:varchar(32);index" json:"status_after"`
	DedupeKey    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"dedupe_key"`
}

func (SystemNotification) TableName() string {
	return "system_notifications"
}

func CreateExecutionLog(db *gorm.DB, item *ExecutionLog) error {
	return db.Create(item).Error
}

func UpsertResourceHealth(db *gorm.DB, item *ResourceHealth) error {
	return db.Save(item).Error
}

func GetResourceHealth(db *gorm.DB, resourceType string, resourceID string) (*ResourceHealth, error) {
	var item ResourceHealth
	if err := db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ListExecutionLogs(db *gorm.DB, query *gorm.DB) ([]*ExecutionLog, error) {
	var items []*ExecutionLog
	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func ListResourceHealth(db *gorm.DB, query *gorm.DB) ([]*ResourceHealth, error) {
	var items []*ResourceHealth
	if err := query.Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func CreateSystemNotification(db *gorm.DB, item *SystemNotification) error {
	return db.Create(item).Error
}

func ListSystemNotifications(db *gorm.DB, query *gorm.DB) ([]*SystemNotification, error) {
	var items []*SystemNotification
	if err := query.Order("created_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
