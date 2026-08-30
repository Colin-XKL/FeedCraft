package util

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCachedFuncWithStore_CoalescesConcurrentMisses(t *testing.T) {
	store := sync.Map{}
	get := func(key string) (string, error) {
		if v, ok := store.Load(key); ok {
			return v.(string), nil
		}
		return "", errors.New("cache miss")
	}
	set := func(value string) error {
		store.Store("stampede-key", value)
		return nil
	}

	var calls int32
	started := make(chan struct{})
	release := make(chan struct{})
	valFunc := func() (string, error) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			close(started)
			<-release
		}
		return "computed", nil
	}

	const workers = 8
	errCh := make(chan error, workers)
	resultCh := make(chan string, workers)
	for i := 0; i < workers; i++ {
		go func() {
			v, err := cachedFuncWithStore("stampede-key", get, set, valFunc, nil)
			errCh <- err
			resultCh <- v
		}()
	}

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("valFunc was not started")
	}
	close(release)

	for i := 0; i < workers; i++ {
		require.NoError(t, <-errCh)
		require.Equal(t, "computed", <-resultCh)
	}
	require.Equal(t, int32(1), atomic.LoadInt32(&calls), "并发缓存未命中应只计算一次")
}
