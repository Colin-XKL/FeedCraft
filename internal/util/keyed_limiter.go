package util

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type semEntry struct {
	sem *semaphore.Weighted
	ref int
}

// KeyedLimiter controls concurrency based on a string key (e.g., domain).
type KeyedLimiter struct {
	limit int64
	mu    sync.Mutex
	sems  map[string]*semEntry // stores semaphores and reference counts
}

// NewKeyedLimiter creates a new KeyedLimiter with the specified default limit.
func NewKeyedLimiter(limit int) *KeyedLimiter {
	return &KeyedLimiter{
		limit: int64(limit),
		sems:  make(map[string]*semEntry),
	}
}

// Acquire blocks until it can acquire a permit for the given key,
// or until the context is canceled. It returns a release function
// that MUST be called to return the permit.
func (l *KeyedLimiter) Acquire(ctx context.Context, key string) (func(), error) {
	l.mu.Lock()
	entry, ok := l.sems[key]
	if !ok {
		entry = &semEntry{
			sem: semaphore.NewWeighted(l.limit),
			ref: 0,
		}
		l.sems[key] = entry
	}
	entry.ref++
	l.mu.Unlock()

	if err := entry.sem.Acquire(ctx, 1); err != nil {
		l.mu.Lock()
		entry.ref--
		if entry.ref == 0 {
			delete(l.sems, key)
		}
		l.mu.Unlock()
		return nil, err
	}

	return func() {
		entry.sem.Release(1)
		l.mu.Lock()
		entry.ref--
		if entry.ref == 0 {
			delete(l.sems, key)
		}
		l.mu.Unlock()
	}, nil
}
