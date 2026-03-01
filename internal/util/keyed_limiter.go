package util

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// KeyedLimiter controls concurrency based on a string key (e.g., domain).
type KeyedLimiter struct {
	limit int64
	sems  sync.Map // stores map[string]*semaphore.Weighted
}

// NewKeyedLimiter creates a new KeyedLimiter with the specified default limit.
func NewKeyedLimiter(limit int) *KeyedLimiter {
	return &KeyedLimiter{limit: int64(limit)}
}

// Acquire blocks until it can acquire a permit for the given key,
// or until the context is canceled. It returns a release function
// that MUST be called to return the permit.
func (l *KeyedLimiter) Acquire(ctx context.Context, key string) (func(), error) {
	v, _ := l.sems.LoadOrStore(key, semaphore.NewWeighted(l.limit))
	sem := v.(*semaphore.Weighted)

	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	return func() { sem.Release(1) }, nil
}
