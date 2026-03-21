package util

import (
	"context"
	"fmt"
	"time"
)

// PriorityDispatcher is a generic worker pool that limits concurrency and supports prioritizing urgent tasks.
// It acts as a global funnel: no matter how many tasks are submitted concurrently,
// only a fixed number of workers will execute them, protecting downstream services from being overwhelmed.
type PriorityDispatcher[R any] struct {
	normalQueue chan taskWrapper[R]
	urgentQueue chan taskWrapper[R]
	// MaxTaskDuration is an optional fallback timeout for each task.
	// If a task takes longer than this, it will be cancelled.
	MaxTaskDuration time.Duration
}

type taskWrapper[R any] struct {
	ctx        context.Context
	fn         func(ctx context.Context) (R, error)
	resultChan chan taskResult[R]
}

type taskResult[R any] struct {
	val R
	err error
}

// NewPriorityDispatcher creates and starts a new PriorityDispatcher with the given number of workers.
func NewPriorityDispatcher[R any](maxConcurrency int) *PriorityDispatcher[R] {
	if maxConcurrency <= 0 {
		panic("priority dispatcher: maxConcurrency must be positive")
	}

	d := &PriorityDispatcher[R]{
		normalQueue: make(chan taskWrapper[R], 1000), // Buffer for pending tasks
		urgentQueue: make(chan taskWrapper[R], 1000),
	}

	// Start fixed number of workers
	for i := 0; i < maxConcurrency; i++ {
		go d.worker()
	}

	return d
}

func (d *PriorityDispatcher[R]) worker() {
	for {
		// ==========================================
		// 优先级调度的核心逻辑 (双层 select 模式)
		// ==========================================
		// 为什么 need 两层 select？
		// 因为在 Go 中如果一层 select 同时有多个 case 就绪，它是随机选择的。
		// 为了实现严格的优先级(重试必须插队)，必须使用两层配合 default 分支。

		// 1. 嗅探层 (非阻塞)：强制优先检查高优队列。
		// 如果 urgentQueue 有任务，立即执行并进入下一轮循环，实现绝对插队。
		// 如果 urgentQueue 为空，则瞬间掉入 default 分支，绝不阻塞等待。
		select {
		case task := <-d.urgentQueue:
			d.executeTask(task)

		default:
			// 2. 阻塞等待层：当确认高优队列为空时，才同时监听两个队列。
			// 如果两个队列都为空，Worker 会在此处被挂起 (Block)，进入零 CPU 消耗的睡眠状态。
			// 直到任一队列有新任务到来，Go 调度器会自动唤醒它。
			select {
			case task := <-d.urgentQueue:
				d.executeTask(task)
			case task := <-d.normalQueue:
				d.executeTask(task)
			}
		}
	}
}

func (d *PriorityDispatcher[R]) executeTask(task taskWrapper[R]) {
	// Panic protection: ensure worker doesn't die and caller gets an error
	defer func() {
		if r := recover(); r != nil {
			var zero R
			// Send the panic as an error to the result channel
			// We use fmt.Errorf for simplicity, but could include stack trace if needed
			task.resultChan <- taskResult[R]{
				val: zero,
				err: fmt.Errorf("task panicked: %v", r),
			}
		}
	}()

	ctx := task.ctx
	if d.MaxTaskDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d.MaxTaskDuration)
		defer cancel()
	}

	// Execute the task with the (possibly timed out) context.
	val, err := task.fn(ctx)
	task.resultChan <- taskResult[R]{val: val, err: err}
}

// Execute submits a task and blocks until it completes or the context is cancelled.
// If urgent is true, the task is sent to the urgentQueue and will be executed before normal tasks.
func (d *PriorityDispatcher[R]) Execute(ctx context.Context, urgent bool, fn func(ctx context.Context) (R, error)) (R, error) {
	resultChan := make(chan taskResult[R], 1)
	task := taskWrapper[R]{
		ctx:        ctx,
		fn:         fn,
		resultChan: resultChan,
	}

	var queue chan taskWrapper[R]
	if urgent {
		queue = d.urgentQueue
	} else {
		queue = d.normalQueue
	}

	// Step 1: Enqueue the task, but respect context cancellation while waiting for queue space.
	select {
	case queue <- task:
		// Task successfully enqueued
	case <-ctx.Done():
		var zero R
		return zero, ctx.Err()
	}

	// Step 2: Wait for the result, but respect context cancellation.
	// Note: Even if we return early due to context cancellation, the worker will still eventually
	// pick up the task and run it (and send the result to resultChan).
	// Because resultChan is buffered (size 1), the worker won't block when sending to it.
	select {
	case res := <-resultChan:
		return res.val, res.err
	case <-ctx.Done():
		var zero R
		return zero, ctx.Err()
	}
}
