package util

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPreheatingScheduler_DefaultQueueSizeMatchesConcurrency(t *testing.T) {
	t.Setenv("FC_PREHEATING_QUEUE_SIZE", "")

	scheduler := NewPreheatingScheduler(func(context.Context, string) error {
		return nil
	}, nil, nil, WithPreheatingMaxConcurrency(3))

	if got := cap(scheduler.taskQueue); got != 3 {
		t.Fatalf("expected default queue size to match concurrency 3, got %d", got)
	}
}

func TestPreheatingScheduler_LimitsConcurrencyAndQueues(t *testing.T) {
	const taskCount = 6
	var running int32
	var maxRunning int32
	var completed int32
	done := make(chan struct{})

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
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

func TestPreheatingScheduler_FullQueueDropsNewestTask(t *testing.T) {
	blockerStarted := make(chan struct{})
	releaseBlocker := make(chan struct{})
	queuedCompleted := make(chan struct{})
	var releaseOnce sync.Once
	var droppedRuns int32
	defer releaseOnce.Do(func() {
		close(releaseBlocker)
	})

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
		switch payload {
		case "blocker":
			close(blockerStarted)
			<-releaseBlocker
		case "queued":
			close(queuedCompleted)
		case "dropped":
			atomic.AddInt32(&droppedRuns, 1)
		}
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(1),
		WithPreheatingInterval(time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	scheduler.ScheduleTask("blocker")
	select {
	case <-blockerStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for blocker preheating task to start")
	}

	scheduler.ScheduleTask("queued")
	waitForPreheatingCondition(t, func() bool {
		return len(scheduler.taskQueue) == 1
	}, "timed out waiting for task to fill preheating queue")

	scheduler.ScheduleTask("dropped")
	waitForPreheatingCondition(t, func() bool {
		return !scheduler.GetContextInfo("dropped").IsActive
	}, "timed out waiting for full queue to reject newest task")

	if got := len(scheduler.taskQueue); got != 1 {
		t.Fatalf("expected accepted queued task to remain in queue, got queue length %d", got)
	}

	releaseOnce.Do(func() {
		close(releaseBlocker)
	})
	select {
	case <-queuedCompleted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for previously queued task to execute")
	}
	if got := atomic.LoadInt32(&droppedRuns); got != 0 {
		t.Fatalf("expected rejected task not to execute, got %d runs", got)
	}
}

func TestPreheatingScheduler_AcceptedTasksRemainFIFO(t *testing.T) {
	blockerStarted := make(chan struct{})
	releaseBlocker := make(chan struct{})
	order := make(chan string, 2)
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() {
		close(releaseBlocker)
	})

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
		if payload == "blocker" {
			close(blockerStarted)
			<-releaseBlocker
			return nil
		}
		order <- payload
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(2),
		WithPreheatingInterval(time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	scheduler.ScheduleTask("blocker")
	select {
	case <-blockerStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for blocker preheating task to start")
	}

	scheduler.ScheduleTask("first")
	waitForPreheatingCondition(t, func() bool {
		return len(scheduler.taskQueue) == 1
	}, "timed out waiting for first queued task")
	scheduler.ScheduleTask("second")
	waitForPreheatingCondition(t, func() bool {
		return len(scheduler.taskQueue) == 2
	}, "timed out waiting for second queued task")

	releaseOnce.Do(func() {
		close(releaseBlocker)
	})
	for _, expected := range []string{"first", "second"} {
		select {
		case actual := <-order:
			if actual != expected {
				t.Fatalf("expected FIFO task %q, got %q", expected, actual)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for FIFO task %q", expected)
		}
	}
}

func TestPreheatingScheduler_DiscardQueuedTaskPreservesNewerVersion(t *testing.T) {
	scheduler := &PreheatingScheduler{
		contexts: map[string]*PreheatingContext{
			"recipe": {
				taskKey:      "recipe",
				timerVersion: 2,
				queued:       true,
			},
		},
	}

	scheduler.discardQueuedTask("recipe", 1)

	if !scheduler.GetContextInfo("recipe").IsActive {
		t.Fatal("expected cleanup for stale timer version to preserve newer context")
	}
}

func TestPreheatingScheduler_DeduplicatesSameRecipeWhileRunning(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	var doneOnce sync.Once
	var runCount int32

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
		if atomic.AddInt32(&runCount, 1) == 1 {
			close(started)
		}
		<-release
		doneOnce.Do(func() {
			close(done)
		})
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(10),
		WithPreheatingInterval(50*time.Millisecond),
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

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
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

func TestPreheatingScheduler_UserRequestInvalidatesQueuedTask(t *testing.T) {
	blockerStarted := make(chan struct{})
	releaseBlocker := make(chan struct{})
	var targetRuns int32

	scheduler := NewPreheatingScheduler(func(_ context.Context, payload string) error {
		if payload == "blocker" {
			close(blockerStarted)
			<-releaseBlocker
			return nil
		}
		atomic.AddInt32(&targetRuns, 1)
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(10),
		WithPreheatingInterval(20*time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
	)

	scheduler.ScheduleTask("blocker")
	select {
	case <-blockerStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for blocker preheating task to start")
	}

	scheduler.ScheduleTask("target")
	time.Sleep(40 * time.Millisecond)

	scheduler.ScheduleTask("target")
	close(releaseBlocker)

	time.Sleep(10 * time.Millisecond)
	if got := atomic.LoadInt32(&targetRuns); got != 0 {
		t.Fatalf("expected stale queued target task to be invalidated, got %d runs", got)
	}

	deadline := time.After(time.Second)
	for atomic.LoadInt32(&targetRuns) == 0 {
		select {
		case <-deadline:
			t.Fatal("timed out waiting for refreshed target task to run")
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func TestPreheatingScheduler_TaskTimeoutReleasesWorker(t *testing.T) {
	slowStarted := make(chan struct{})
	slowCancelled := make(chan struct{})
	fastCompleted := make(chan struct{})

	scheduler := NewPreheatingScheduler(func(ctx context.Context, payload string) error {
		if payload == "slow" {
			close(slowStarted)
			<-ctx.Done()
			close(slowCancelled)
			return ctx.Err()
		}
		close(fastCompleted)
		return nil
	}, nil, nil,
		WithPreheatingMaxConcurrency(1),
		WithPreheatingQueueSize(10),
		WithPreheatingInterval(time.Millisecond),
		WithPreheatingJitter(0),
		WithPreheatingMaxCount(1),
		WithPreheatingGraceTime(time.Hour),
		WithPreheatingTaskTimeout(20*time.Millisecond),
	)

	scheduler.ScheduleTask("slow")
	select {
	case <-slowStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for slow task to start")
	}

	scheduler.ScheduleTask("fast")
	select {
	case <-slowCancelled:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for slow task context cancellation")
	}
	select {
	case <-fastCompleted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for worker to execute queued task after cancellation")
	}
}

func waitForPreheatingCondition(t *testing.T, condition func() bool, failureMessage string) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal(failureMessage)
}
