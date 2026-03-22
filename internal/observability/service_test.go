package observability

import (
	"FeedCraft/internal/dao"
	"testing"
	"time"
)

func TestApplyEventToHealthTransitions(t *testing.T) {
	now := time.Now()
	health := &dao.ResourceHealth{
		ResourceType:  ResourceTypeRecipe,
		ResourceID:    "demo",
		ResourceName:  "demo",
		CurrentStatus: dao.ResourceStatusHealthy,
	}

	for i := 0; i < DegradedFailureThreshold; i++ {
		applyEventToHealth(health, ExecutionEvent{
			ResourceType: ResourceTypeRecipe,
			ResourceID:   "demo",
			Status:       dao.ExecutionStatusFailure,
			ErrorKind:    ErrorKindNetwork,
			Message:      "network failed",
			OccurredAt:   now.Add(time.Duration(i) * time.Minute),
		})
	}

	if health.CurrentStatus != dao.ResourceStatusDegraded {
		t.Fatalf("expected degraded after %d failures, got %s", DegradedFailureThreshold, health.CurrentStatus)
	}

	for i := DegradedFailureThreshold; i < PausedFailureThreshold; i++ {
		applyEventToHealth(health, ExecutionEvent{
			ResourceType: ResourceTypeRecipe,
			ResourceID:   "demo",
			Status:       dao.ExecutionStatusFailure,
			ErrorKind:    ErrorKindNetwork,
			Message:      "network failed",
			OccurredAt:   now.Add(time.Duration(i) * time.Minute),
		})
	}

	if health.CurrentStatus != dao.ResourceStatusPaused {
		t.Fatalf("expected paused after %d failures, got %s", PausedFailureThreshold, health.CurrentStatus)
	}
	if health.PausedAt == nil {
		t.Fatal("expected paused_at to be set")
	}
}

func TestBuildNotificationOnStatusTransition(t *testing.T) {
	now := time.Now()
	health := &dao.ResourceHealth{
		ResourceType:  ResourceTypeTopic,
		ResourceID:    "topic-1",
		ResourceName:  "Topic 1",
		CurrentStatus: dao.ResourceStatusDegraded,
	}
	item := buildNotification(dao.ResourceStatusHealthy, health, ExecutionEvent{
		ResourceType: ResourceTypeTopic,
		ResourceID:   "topic-1",
		Status:       dao.ExecutionStatusPartialSuccess,
		Message:      "partial upstream failures",
		OccurredAt:   now,
	})
	if item == nil {
		t.Fatal("expected a notification for degraded transition")
	}
	if item.EventType != "resource_degraded" {
		t.Fatalf("unexpected event type: %s", item.EventType)
	}
}
