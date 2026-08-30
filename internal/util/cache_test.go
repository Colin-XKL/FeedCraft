package util

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func setupTestRedis(tb testing.TB) *miniredis.Miniredis {
	tb.Helper()

	server, err := miniredis.Run()
	if err != nil {
		tb.Fatalf("could not start miniredis: %v", err)
	}
	tb.Setenv("FC_REDIS_URI", fmt.Sprintf("redis://%s", server.Addr()))
	tb.Cleanup(server.Close)
	return server
}

func TestCachedFuncCoalescesConcurrentMisses(t *testing.T) {
	setupTestRedis(t)

	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})

	valFunc := func() (string, error) {
		if calls.Add(1) == 1 {
			close(started)
		}
		<-release
		return "computed-once", nil
	}

	const goroutines = 8
	errCh := make(chan error, goroutines)
	results := make([]string, goroutines)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			value, err := CachedFunc("stampede-key", valFunc)
			if err != nil {
				errCh <- err
				return
			}
			results[idx] = value
		}(i)
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("valFunc was never invoked")
	}
	time.Sleep(50 * time.Millisecond)
	close(release)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("CachedFunc returned error: %v", err)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("expected valFunc to run once for concurrent cache misses, got %d calls", got)
	}
	for i, value := range results {
		if value != "computed-once" {
			t.Fatalf("result[%d] = %q, want computed-once", i, value)
		}
	}
}

func TestCachedFuncDifferentKeysDoNotShareFlight(t *testing.T) {
	setupTestRedis(t)

	var calls atomic.Int32
	valFunc := func() (string, error) {
		n := calls.Add(1)
		return fmt.Sprintf("value-%d", n), nil
	}

	first, err := CachedFunc("key-a", valFunc)
	if err != nil {
		t.Fatalf("key-a: %v", err)
	}
	second, err := CachedFunc("key-b", valFunc)
	if err != nil {
		t.Fatalf("key-b: %v", err)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected 2 valFunc calls for different keys, got %d", calls.Load())
	}
	if first == second {
		t.Fatalf("different keys unexpectedly shared result %q", first)
	}
}

func TestCachedFuncHitSkipsValFunc(t *testing.T) {
	setupTestRedis(t)

	var calls atomic.Int32
	valFunc := func() (string, error) {
		calls.Add(1)
		return "fresh", nil
	}

	if _, err := CachedFunc("hit-key", valFunc); err != nil {
		t.Fatalf("first fill: %v", err)
	}
	value, err := CachedFunc("hit-key", valFunc)
	if err != nil {
		t.Fatalf("second read: %v", err)
	}
	if value != "fresh" {
		t.Fatalf("cached value = %q, want fresh", value)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected valFunc once on cache hit path, got %d", calls.Load())
	}
}

func TestCachedFuncErrorIsNotCached(t *testing.T) {
	setupTestRedis(t)

	var calls atomic.Int32
	valFunc := func() (string, error) {
		if calls.Add(1) == 1 {
			return "", errors.New("llm timeout")
		}
		return "recovered", nil
	}

	if _, err := CachedFunc("err-key", valFunc); err == nil {
		t.Fatal("expected first call to fail")
	}
	value, err := CachedFunc("err-key", valFunc)
	if err != nil {
		t.Fatalf("retry after error: %v", err)
	}
	if value != "recovered" {
		t.Fatalf("retry value = %q, want recovered", value)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected valFunc to retry after error, got %d calls", calls.Load())
	}
}
