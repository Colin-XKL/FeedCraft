package dao

import "testing"

func TestTopicFeed_NormalizeInputs_FromInputs(t *testing.T) {
	topic := &TopicFeed{
		Inputs: []TopicInput{
			{URI: "https://example.com/a.xml", Description: "Tech news"},
			{URI: " feedcraft://recipe/foo ", Description: ""},
		},
	}

	topic.NormalizeInputs()

	if len(topic.InputURIs) != 2 {
		t.Fatalf("expected 2 input URIs, got %d", len(topic.InputURIs))
	}
	if topic.InputURIs[0] != "https://example.com/a.xml" {
		t.Fatalf("unexpected first URI: %q", topic.InputURIs[0])
	}
	if topic.InputURIs[1] != "feedcraft://recipe/foo" {
		t.Fatalf("unexpected second URI: %q", topic.InputURIs[1])
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

	if len(topic.InputURIs) != 1 {
		t.Fatalf("expected 1 enabled URI, got %d", len(topic.InputURIs))
	}
	if topic.InputURIs[0] != "https://example.com/enabled.xml" {
		t.Fatalf("unexpected enabled URI: %q", topic.InputURIs[0])
	}
	if len(topic.Inputs) != 2 {
		t.Fatalf("expected 2 inputs retained, got %d", len(topic.Inputs))
	}
	if !topic.Inputs[1].Disabled {
		t.Fatal("expected second input to stay disabled")
	}
}

func TestTopicFeed_NormalizeInputs_FromLegacyURIs(t *testing.T) {
	topic := &TopicFeed{
		InputURIs: []string{"https://example.com/a.xml", "https://example.com/b.xml"},
	}

	topic.NormalizeInputs()

	if len(topic.Inputs) != 2 {
		t.Fatalf("expected 2 inputs, got %d", len(topic.Inputs))
	}
	if topic.Inputs[0].URI != "https://example.com/a.xml" {
		t.Fatalf("unexpected first input URI: %q", topic.Inputs[0].URI)
	}
}
