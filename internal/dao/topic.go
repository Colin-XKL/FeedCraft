package dao

import (
	"strings"

	"gorm.io/gorm"
)

// TopicInput is a single upstream source with optional admin metadata.
type TopicInput struct {
	URI         string `json:"uri"`
	Description string `json:"description,omitempty"`
	// Disabled excludes this input from topic aggregation when true.
	Disabled bool `json:"disabled,omitempty"`
}

// TopicFeed represents the persistence model for a multi-source aggregation node.
type TopicFeed struct {
	BaseModelWithoutPK
	ID          string `gorm:"primaryKey" json:"id" binding:"required"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	// Inputs carries URI plus optional description for admin display.
	Inputs []TopicInput `json:"inputs,omitempty" gorm:"serializer:json"`

	// List of URIs representing inputs (derived from Inputs for runtime).
	// Uses a custom protocol for internal resources to make routing elegant and standard.
	// Examples:
	//   - "feedcraft://recipe/my-tech-recipe" (Internal RecipeFeed)
	//   - "feedcraft://topic/sub-topic-id"    (Nested internal TopicFeed)
	//   - "https://external.com/rss.xml"      (External raw feed)
	InputURIs []string `json:"input_uris" gorm:"serializer:json"`

	// Configuration for the aggregator pipeline
	AggregatorConfig []AggregatorStep `json:"aggregator_config" gorm:"serializer:json"`
}

// NormalizeInputs keeps Inputs and InputURIs in sync.
// When Inputs is provided it becomes the source of truth; otherwise legacy InputURIs are upgraded.
func (t *TopicFeed) NormalizeInputs() {
	if t == nil {
		return
	}

	if len(t.Inputs) > 0 {
		uris := make([]string, 0, len(t.Inputs))
		normalized := make([]TopicInput, 0, len(t.Inputs))
		for _, item := range t.Inputs {
			uri := strings.TrimSpace(item.URI)
			if uri == "" {
				continue
			}
			if !item.Disabled {
				uris = append(uris, uri)
			}
			normalized = append(normalized, TopicInput{
				URI:         uri,
				Description: strings.TrimSpace(item.Description),
				Disabled:    item.Disabled,
			})
		}
		t.InputURIs = uris
		t.Inputs = normalized
		return
	}

	if len(t.InputURIs) == 0 {
		t.Inputs = nil
		return
	}

	inputs := make([]TopicInput, 0, len(t.InputURIs))
	for _, uri := range t.InputURIs {
		uri = strings.TrimSpace(uri)
		if uri == "" {
			continue
		}
		inputs = append(inputs, TopicInput{URI: uri})
	}
	t.Inputs = inputs
	t.InputURIs = make([]string, 0, len(inputs))
	for _, item := range inputs {
		t.InputURIs = append(t.InputURIs, item.URI)
	}
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
	if result.Error != nil {
		return nil, result.Error
	}
	return &topic, nil
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
		return nil, err
	}
	return topics, nil
}
