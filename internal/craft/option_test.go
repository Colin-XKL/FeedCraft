package craft

import (
	"errors"
	"testing"

	"github.com/gorilla/feeds"
)

func TestOptionTransformFeedItem_PartialFailure(t *testing.T) {
	// Setup
	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Item 1"}, // Will fail
			{Title: "Item 2"}, // Will succeed
		},
	}
	payload := ExtraPayload{}

	// Processor that fails on the first item
	processor := func(item *feeds.Item, p ExtraPayload) error {
		if item.Title == "Item 1" {
			return errors.New("failed item 1")
		}
		return nil
	}

	// Execution
	option := OptionTransformFeedItem(processor)
	err := option(feed, payload)

	// Verify - New behavior: should NOT fail because at least one item succeeded
	if err != nil {
		t.Errorf("Expected no error for partial failure, got: %v", err)
	}
}

func TestOptionTransformFeedItem_AllFailure(t *testing.T) {
	// Setup
	feed := &feeds.Feed{
		Items: []*feeds.Item{
			{Title: "Item 1"},
			{Title: "Item 2"},
		},
	}
	payload := ExtraPayload{}

	// Processor that always fails
	processor := func(item *feeds.Item, p ExtraPayload) error {
		return errors.New("failed")
	}

	// Execution
	option := OptionTransformFeedItem(processor)
	err := option(feed, payload)

	// Verify - New behavior: should fail because ALL items failed
	if err == nil {
		t.Errorf("Expected error when all items fail, got nil")
	}
}
