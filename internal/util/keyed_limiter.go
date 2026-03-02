package util

import (
	"context"
	"hash/fnv"

	"golang.org/x/sync/semaphore"
)

const defaultBucketSize = 256

// KeyedLimiter controls concurrency based on a string key (e.g., domain)
// using a fixed-size hash bucket strategy to ensure OOM safety and constant memory.
type KeyedLimiter struct {
	limit   int64
	buckets []*semaphore.Weighted
	size    uint32
}

// NewKeyedLimiter creates a new KeyedLimiter with the specified default limit.
// It uses a fixed number of buckets (256) to balance collision risk and memory usage.
func NewKeyedLimiter(limit int) *KeyedLimiter {
	l := &KeyedLimiter{
		limit:   int64(limit),
		size:    defaultBucketSize,
		buckets: make([]*semaphore.Weighted, defaultBucketSize),
	}
	for i := range l.buckets {
		l.buckets[i] = semaphore.NewWeighted(l.limit)
	}
	return l
}

// Acquire blocks until it can acquire a permit for the given key,
// or until the context is canceled. It returns a release function
// that MUST be called to return the permit.
func (l *KeyedLimiter) Acquire(ctx context.Context, key string) (func(), error) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := h.Sum32() % l.size

	sem := l.buckets[idx]

	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	return func() { sem.Release(1) }, nil
}
