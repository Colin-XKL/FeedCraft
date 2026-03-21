package util

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestPriorityDispatcher_PanicRecovery(t *testing.T) {
	// 1 worker to ensure we test if it survives the panic
	dispatcher := NewPriorityDispatcher[string](1)

	// Task 1: Trigger a panic
	t.Log("Starting panicking task...")
	_, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
		panic("boom")
	})

	if err == nil {
		t.Fatal("expected error from panicked task, got nil")
	}
	if !strings.Contains(err.Error(), "task panicked: boom") {
		t.Errorf("expected error message to contain 'task panicked: boom', got: %v", err)
	}

	// Task 2: Check if worker is still alive and can process new tasks
	t.Log("Starting recovery verification task...")
	val, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
		return "alive", nil
	})

	if err != nil {
		t.Fatalf("worker died after panic: %v", err)
	}
	if val != "alive" {
		t.Errorf("expected 'alive', got: %s", val)
	}
}

func TestPriorityDispatcher_Timeout(t *testing.T) {
	dispatcher := NewPriorityDispatcher[string](1)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := dispatcher.Execute(ctx, false, func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "done", nil
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestPriorityDispatcher_MaxTaskDuration(t *testing.T) {
	dispatcher := NewPriorityDispatcher[string](1)
	dispatcher.MaxTaskDuration = 50 * time.Millisecond

	start := time.Now()
	_, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
		<-ctx.Done() // Wait for dispatcher to cancel
		return "", ctx.Err()
	})

	duration := time.Since(start)
	if duration > 100*time.Millisecond {
		t.Errorf("task took too long: %v", duration)
	}
	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("expected timeout error from MaxTaskDuration, got: %v", err)
	}
}

func TestPriorityDispatcher_Priority(t *testing.T) {
	// 1 worker, so tasks run sequentially
	dispatcher := NewPriorityDispatcher[string](1)

	// Block the worker with a long task
	blockChan := make(chan struct{})
	workerStarted := make(chan struct{})
	errChan := make(chan error, 4)

	go func() {
		_, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
			close(workerStarted)
			<-blockChan
			return "unblocked", nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// Wait to ensure the blocking task is picked up
	<-workerStarted

	// Enqueue 2 normal tasks
	order := make(chan string, 3)
	go func() {
		_, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
			order <- "normal 1"
			return "", nil
		})
		if err != nil {
			errChan <- err
		}
	}()
	go func() {
		_, err := dispatcher.Execute(context.Background(), false, func(ctx context.Context) (string, error) {
			order <- "normal 2"
			return "", nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// Wait to ensure normal tasks are in queue
	for len(dispatcher.normalQueue) < 2 {
		time.Sleep(1 * time.Millisecond)
	}

	// Enqueue 1 urgent task
	go func() {
		_, err := dispatcher.Execute(context.Background(), true, func(ctx context.Context) (string, error) {
			order <- "urgent"
			return "", nil
		})
		if err != nil {
			errChan <- err
		}
	}()

	// Wait to ensure urgent task is in queue
	for len(dispatcher.urgentQueue) < 1 {
		time.Sleep(1 * time.Millisecond)
	}

	// Unblock the worker
	close(blockChan)

	// The first task out should be "urgent" even though it was enqueued last
	select {
	case first := <-order:
		if first != "urgent" {
			t.Errorf("expected urgent task to be first after unblocking, got: %s", first)
		}
	case err := <-errChan:
		t.Fatalf("unexpected execute error: %v", err)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for urgent task")
	}

	// Wait for the remaining normal tasks
	for i := 0; i < 2; i++ {
		select {
		case <-order:
		case err := <-errChan:
			t.Fatalf("unexpected execute error: %v", err)
		case <-time.After(1 * time.Second):
			t.Fatal("timeout waiting for normal tasks")
		}
	}
}
