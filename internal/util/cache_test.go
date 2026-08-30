package util

import (
	"context"
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

func TestCachedFuncWithStoreContext_CanceledFollowerStopsWaiting(t *testing.T) {
	get := func(string) (string, error) { return "", errors.New("cache miss") }
	set := func(string) error { return nil }
	started := make(chan struct{})
	release := make(chan struct{})
	valFunc := func(context.Context) (string, error) {
		close(started)
		<-release
		return "computed", nil
	}

	leaderDone := make(chan error, 1)
	go func() {
		_, err := cachedFuncWithStoreContext(context.Background(), "follower-cancel", get, set, valFunc, nil)
		leaderDone <- err
	}()
	<-started

	followerCtx, cancelFollower := context.WithCancel(context.Background())
	followerDone := make(chan error, 1)
	go func() {
		_, err := cachedFuncWithStoreContext(followerCtx, "follower-cancel", get, set, valFunc, nil)
		followerDone <- err
	}()
	waitForCacheWaiters(t, "follower-cancel", 2)
	cancelFollower()

	select {
	case err := <-followerDone:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("canceled follower did not stop waiting")
	}
	close(release)
	require.NoError(t, <-leaderDone)
}

func TestCachedFuncWithStoreContext_CanceledLeaderDoesNotFailFollower(t *testing.T) {
	get := func(string) (string, error) { return "", errors.New("cache miss") }
	set := func(string) error { return nil }
	started := make(chan struct{})
	release := make(chan struct{})
	var calls int32
	valFunc := func(ctx context.Context) (string, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			close(started)
		}
		select {
		case <-release:
			return "computed", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	leaderDone := make(chan error, 1)
	go func() {
		_, err := cachedFuncWithStoreContext(leaderCtx, "leader-cancel", get, set, valFunc, nil)
		leaderDone <- err
	}()
	<-started

	followerDone := make(chan struct {
		value string
		err   error
	}, 1)
	go func() {
		value, err := cachedFuncWithStoreContext(context.Background(), "leader-cancel", get, set, valFunc, nil)
		followerDone <- struct {
			value string
			err   error
		}{value, err}
	}()
	waitForCacheWaiters(t, "leader-cancel", 2)
	cancelLeader()
	require.ErrorIs(t, <-leaderDone, context.Canceled)
	close(release)

	result := <-followerDone
	require.NoError(t, result.err)
	require.Equal(t, "computed", result.value)
	require.Equal(t, int32(1), atomic.LoadInt32(&calls))
}

func waitForCacheWaiters(t *testing.T, key string, want int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		cacheFlight.mu.Lock()
		call := cacheFlight.calls[key]
		got := 0
		if call != nil {
			got = call.waiters
		}
		cacheFlight.mu.Unlock()
		if got >= want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("cache flight %q did not reach %d waiters", key, want)
}
