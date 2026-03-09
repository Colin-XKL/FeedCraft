package util

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestKeyedLimiter_Basic(t *testing.T) {
	limiter := NewKeyedLimiter(1)

	// Acquire for key1
	release, err := limiter.Acquire(context.Background(), "key1")
	if err != nil {
		t.Fatalf("failed to acquire: %v", err)
	}

	// Try to acquire again for key1 (should block)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err = limiter.Acquire(ctx, "key1")
	if err == nil {
		t.Error("expected timeout when acquiring same key twice with limit 1")
	}

	release()

	// Now should be able to acquire again
	release2, err := limiter.Acquire(context.Background(), "key1")
	if err != nil {
		t.Fatalf("failed to acquire after release: %v", err)
	}
	release2()
}

func TestKeyedLimiter_Concurrency(t *testing.T) {
	limit := 5
	limiter := NewKeyedLimiter(limit)
	key := "domain.com"

	var wg sync.WaitGroup
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rel, err := limiter.Acquire(context.Background(), key)
			if err != nil {
				t.Errorf("failed to acquire: %v", err)
				return
			}
			time.Sleep(10 * time.Millisecond)
			rel()
		}()
	}
	wg.Wait()
}

func TestKeyedLimiter_DifferentKeys(t *testing.T) {
	limiter := NewKeyedLimiter(1)

	// Acquire for key1
	rel1, _ := limiter.Acquire(context.Background(), "a.com")
	defer rel1()

	// Acquire for key2 (should not block, unless collision happens)
	// With 1024 buckets, the chance of collision for 2 keys is 1/1024.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	rel2, err := limiter.Acquire(ctx, "b.com")
	if err != nil {
		t.Errorf("key2 blocked by key1 (possible collision, but unlikely): %v", err)
	} else if rel2 != nil {
		rel2()
	}
}
