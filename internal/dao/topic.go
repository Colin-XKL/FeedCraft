package dao

import (
	"gorm.io/gorm"
)

// TopicFeed represents the persistence model for a multi-source aggregation node.
type TopicFeed struct {
	BaseModelWithoutPK
	ID          string `gorm:"primaryKey" json:"id" binding:"required"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	// List of URIs representing inputs.
	// Uses a custom protocol for internal resources to make routing elegant and standard.
	// Examples:
	//   - "feedcraft://recipe/my-tech-recipe" (Internal RecipeFeed)
	//   - "feedcraft://topic/sub-topic-id"    (Nested internal TopicFeed)
	//   - "https://external.com/rss.xml"      (External raw feed)
	InputURIs []string `json:"input_uris" binding:"required" gorm:"serializer:json"`

	// Configuration for the aggregator pipeline
	AggregatorConfig []AggregatorStep `json:"aggregator_config" gorm:"serializer:json"`
}

// AggregatorStep defines a single processing step in an Aggregator pipeline.
type AggregatorStep struct {
	Type   string            `json:"type" binding:"required"` // e.g., "deduplicate", "sort", "limit"
	Option map[string]string `json:"option"`                  // e.g., {"by": "date_desc"} or {"max": "50"}
}

// TableName overrides the default table name for TopicFeed.
func (TopicFeed) TableName() string {
	return "topic_feeds"
}

// CreateTopicFeed creates a new TopicFeed record in the database.
func CreateTopicFeed(db *gorm.DB, topic *TopicFeed) error {
	return db.Create(topic).Error
}

// GetTopicFeedByID retrieves a TopicFeed record by its ID.
func GetTopicFeedByID(db *gorm.DB, id string) (*TopicFeed, error) {
	var topic TopicFeed
	if id == "" {
		return nil, gorm.ErrRecordNotFound
	}
	result := db.Where("id = ?", id).First(&topic)
	return &topic, result.Error
}

// UpdateTopicFeed updates an existing TopicFeed record.
func UpdateTopicFeed(db *gorm.DB, topic *TopicFeed) error {
	return db.Save(topic).Error
}

// DeleteTopicFeed deletes a TopicFeed record by its ID.
func DeleteTopicFeed(db *gorm.DB, id string) error {
	var topic TopicFeed
	result := db.Where("id = ?", id).Delete(&topic)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListTopicFeeds retrieves all TopicFeed records.
func ListTopicFeeds(db *gorm.DB) ([]*TopicFeed, error) {
	var topics []*TopicFeed
	if err := db.Find(&topics).Error; err != nil {
		return topics, err
	}
	return topics, nil
}
