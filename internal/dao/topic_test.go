package dao

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestTopicFeed_NormalizeInputs_FromInputs(t *testing.T) {
	topic := &TopicFeed{
		Inputs: []TopicInput{
			{URI: "https://example.com/a.xml", Description: "Tech news"},
			{URI: " feedcraft://recipe/foo ", Description: ""},
		},
	}

	topic.NormalizeInputs()

	enabled := topic.EnabledInputURIs()
	if len(enabled) != 2 {
		t.Fatalf("expected 2 enabled input URIs, got %d", len(enabled))
	}
	if enabled[0] != "https://example.com/a.xml" {
		t.Fatalf("unexpected first URI: %q", enabled[0])
	}
	if enabled[1] != "feedcraft://recipe/foo" {
		t.Fatalf("unexpected second URI: %q", enabled[1])
	}
	if topic.Inputs[0].Description != "Tech news" {
		t.Fatalf("unexpected description: %q", topic.Inputs[0].Description)
	}
}

func TestTopicFeed_NormalizeInputs_SkipsDisabled(t *testing.T) {
	topic := &TopicFeed{
		Inputs: []TopicInput{
			{URI: "https://example.com/enabled.xml"},
			{URI: "https://example.com/disabled.xml", Disabled: true},
		},
	}

	topic.NormalizeInputs()

	enabled := topic.EnabledInputURIs()
	if len(enabled) != 1 {
		t.Fatalf("expected 1 enabled URI, got %d", len(enabled))
	}
	if enabled[0] != "https://example.com/enabled.xml" {
		t.Fatalf("unexpected enabled URI: %q", enabled[0])
	}
	if len(topic.Inputs) != 2 {
		t.Fatalf("expected 2 inputs retained, got %d", len(topic.Inputs))
	}
	if !topic.Inputs[1].Disabled {
		t.Fatal("expected second input to stay disabled")
	}
}

func TestMigrateTopicFeedInputs_FromLegacyInputURIs(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE topic_feeds (
			id TEXT PRIMARY KEY,
			title TEXT,
			input_uris TEXT,
			aggregator_config TEXT,
			created_at DATETIME,
			updated_at DATETIME
		)
	`).Error; err != nil {
		t.Fatalf("create legacy topic_feeds table: %v", err)
	}
	if err := db.Exec(
		`INSERT INTO topic_feeds (id, title, input_uris, aggregator_config) VALUES (?, ?, ?, ?)`,
		"legacy-topic",
		"Legacy Topic",
		`["https://example.com/a.xml","feedcraft://recipe/foo"]`,
		`[]`,
	).Error; err != nil {
		t.Fatalf("insert legacy topic feed: %v", err)
	}

	if err := db.AutoMigrate(&TopicFeed{}); err != nil {
		t.Fatalf("auto migrate topic feed: %v", err)
	}
	migrateTopicFeedInputs(db)

	var topic TopicFeed
	if err := db.First(&topic, "id = ?", "legacy-topic").Error; err != nil {
		t.Fatalf("load migrated topic feed: %v", err)
	}
	if len(topic.Inputs) != 2 {
		t.Fatalf("expected 2 migrated inputs, got %d", len(topic.Inputs))
	}
	if topic.Inputs[0].URI != "https://example.com/a.xml" {
		t.Fatalf("unexpected first migrated input URI: %q", topic.Inputs[0].URI)
	}
	if db.Migrator().HasColumn("topic_feeds", "input_uris") {
		t.Fatal("expected legacy input_uris column to be dropped")
	}
}
