package util

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPreheatingScheduler_LimitsConcurrencyAndQueues(t *testing.T) {
	const taskCount = 6
	var running int32
	var maxRunning int32
	var completed int32
	done := make(chan struct{})

	scheduler := NewPreheatingScheduler(func(payload string) error {
		current := atomic.AddInt32(&running, 1)
		for {
			previous := atomic.LoadInt32(&maxRunning)
			if current <= previous || atomic.CompareAndSwapInt32(&maxRunning, previous, current) {
				break
			}
		}
		time.Sleep(20 * time.Millisecond)
		atomic.AddInt32(&running, -1)
		if atomic.AddInt32(&completed, 1) == taskCount {
			close(done)
		}
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(2),
		WithPreheatingQueueSize(taskCount),
		WithPreheatingInterval(time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	for i := 0; i < taskCount; i++ {
		scheduler.ScheduleTask(string(rune('a' + i)))
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for queued preheating tasks, completed=%d", atomic.LoadInt32(&completed))
	}

	if got := atomic.LoadInt32(&maxRunning); got > 2 {
		t.Fatalf("expected at most 2 concurrent tasks, got %d", got)
	}
}

func TestPreheatingScheduler_DeduplicatesSameRecipeWhileRunning(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	var runCount int32

	scheduler := NewPreheatingScheduler(func(payload string) error {
		if atomic.AddInt32(&runCount, 1) == 1 {
			close(started)
		}
		<-release
		close(done)
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(10),
		WithPreheatingInterval(time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	scheduler.ScheduleTask("same")
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first preheating task to start")
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scheduler.ScheduleTask("same")
		}()
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	close(release)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for preheating task to finish")
	}

	if got := atomic.LoadInt32(&runCount); got != 1 {
		t.Fatalf("expected duplicate schedules for the same recipe to be coalesced while running, got %d runs", got)
	}
}

func TestPreheatingScheduler_UserRequestsDoNotConsumePreheatingCount(t *testing.T) {
	done := make(chan struct{})
	var runCount int32

	scheduler := NewPreheatingScheduler(func(payload string) error {
		if atomic.AddInt32(&runCount, 1) == 1 {
			close(done)
		}
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(10),
		WithPreheatingInterval(20*time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	for i := 0; i < 5; i++ {
		scheduler.ScheduleTask("frequent")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("frequent user requests consumed the preheating count before any preheating task ran")
	}
}
