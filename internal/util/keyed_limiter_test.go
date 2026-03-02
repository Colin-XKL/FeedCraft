package util

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestKeyedLimiterEviction(t *testing.T) {
	limiter := NewKeyedLimiter(2)
	ctx := context.Background()

	release, err := limiter.Acquire(ctx, "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// At this point, map should have 1 entry
	limiter.mu.Lock()
	if len(limiter.sems) != 1 {
		t.Errorf("expected 1 entry, got %d", len(limiter.sems))
	}
	limiter.mu.Unlock()

	release()

	// After release, map should have 0 entries
	limiter.mu.Lock()
	if len(limiter.sems) != 0 {
		t.Errorf("expected 0 entries, got %d", len(limiter.sems))
	}
	limiter.mu.Unlock()
}

func TestKeyedLimiterConcurrency(t *testing.T) {
	limiter := NewKeyedLimiter(2)
	ctx := context.Background()
	key := "example.com"

	var wg sync.WaitGroup
	var activeCount int
	var mu sync.Mutex

	// Launch 5 goroutines to acquire the semaphore
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release, err := limiter.Acquire(ctx, key)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			mu.Lock()
			activeCount++
			if activeCount > 2 {
				t.Errorf("exceeded max concurrency of 2, active: %d", activeCount)
			}
			mu.Unlock()

			time.Sleep(50 * time.Millisecond) // Simulate work

			mu.Lock()
			activeCount--
			mu.Unlock()

			release()
		}()
	}

	wg.Wait()

	// After all releases, map should have 0 entries
	limiter.mu.Lock()
	if len(limiter.sems) != 0 {
		t.Errorf("expected 0 entries after all goroutines finish, got %d", len(limiter.sems))
	}
	limiter.mu.Unlock()
}

func TestKeyedLimiterMultipleKeys(t *testing.T) {
	limiter := NewKeyedLimiter(1)
	ctx := context.Background()

	release1, err := limiter.Acquire(ctx, "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	release2, err := limiter.Acquire(ctx, "test.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	limiter.mu.Lock()
	if len(limiter.sems) != 2 {
		t.Errorf("expected 2 entries, got %d", len(limiter.sems))
	}
	limiter.mu.Unlock()

	release1()

	limiter.mu.Lock()
	if len(limiter.sems) != 1 {
		t.Errorf("expected 1 entry, got %d", len(limiter.sems))
	}
	if _, ok := limiter.sems["test.com"]; !ok {
		t.Errorf("expected test.com to still be in map")
	}
	limiter.mu.Unlock()

	release2()

	limiter.mu.Lock()
	if len(limiter.sems) != 0 {
		t.Errorf("expected 0 entries, got %d", len(limiter.sems))
	}
	limiter.mu.Unlock()
}

func TestKeyedLimiterContextCancellation(t *testing.T) {
	limiter := NewKeyedLimiter(1)
	ctx, cancel := context.WithCancel(context.Background())

	release1, err := limiter.Acquire(ctx, "example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create a channel to know when the second acquire fails
	errCh := make(chan error)
	go func() {
		_, err := limiter.Acquire(ctx, "example.com")
		errCh <- err
	}()

	// Ensure the second acquire has started
	time.Sleep(10 * time.Millisecond)

	limiter.mu.Lock()
	if l := len(limiter.sems); l != 1 {
		t.Errorf("expected 1 entry, got %d", l)
	}
	if entry := limiter.sems["example.com"]; entry.ref != 2 {
		t.Errorf("expected ref to be 2, got %d", entry.ref)
	}
	limiter.mu.Unlock()

	// Cancel context to force the blocked acquire to return error
	cancel()
	err = <-errCh
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	// After cancellation, the ref count should be decremented but the entry still exists
	limiter.mu.Lock()
	if l := len(limiter.sems); l != 1 {
		t.Errorf("expected 1 entry, got %d", l)
	}
	if entry := limiter.sems["example.com"]; entry.ref != 1 {
		t.Errorf("expected ref to be 1, got %d", entry.ref)
	}
	limiter.mu.Unlock()

	release1()

	limiter.mu.Lock()
	if l := len(limiter.sems); l != 0 {
		t.Errorf("expected 0 entries, got %d", l)
	}
	limiter.mu.Unlock()
}
