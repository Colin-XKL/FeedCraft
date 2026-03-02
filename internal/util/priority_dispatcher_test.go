package util

import (
	"context"
	"testing"
	"time"
)

func TestPriorityDispatcher_Timeout(t *testing.T) {
	// Create a dispatcher with 1 worker
	d := NewPriorityDispatcher[string](1)

	// Submit a task that blocks forever
	// We use a context that never expires for the first task to keep the worker busy
	go func() {
		_, _ = d.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
			select {} // Block forever
		})
	}()

	// Wait a bit to ensure the worker is picked up the task
	time.Sleep(100 * time.Millisecond)

	// Now try to submit another task with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := d.Execute(ctx, false, func(ctx context.Context) (string, error) {
		return "ok", nil
	})
	duration := time.Since(start)

	if err == nil {
		t.Fatal("Task should have failed with timeout")
	}
	if err != context.DeadlineExceeded {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}

	// Ensure it didn't take too long
	if duration > 500*time.Millisecond {
		t.Fatalf("Timeout took too long: %v", duration)
	}
	t.Logf("Task correctly timed out after %v", duration)
}

func TestPriorityDispatcher_TaskRespectsContext(t *testing.T) {
	d := NewPriorityDispatcher[string](1)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := d.Execute(ctx, false, func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(500 * time.Millisecond):
			return "too late", nil
		}
	})

	if err != context.DeadlineExceeded {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}
}

func TestPriorityDispatcher_MaxTaskDuration(t *testing.T) {
	d := NewPriorityDispatcher[string](1)
	d.MaxTaskDuration = 100 * time.Millisecond

	// Even if we use context.Background(), the dispatcher should time out the task
	_, err := d.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(1 * time.Second):
			return "too late", nil
		}
	})

	if err != context.DeadlineExceeded {
		t.Fatalf("Expected DeadlineExceeded from dispatcher, got %v", err)
	}
}
