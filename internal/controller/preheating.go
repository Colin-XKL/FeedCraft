package controller

import (
	"log"
	"sync"
	"time"
)

// PreheatingScheduler 预热调度器
type PreheatingScheduler struct {
	tasks    map[string]*time.Timer // key: taskID, value: timer
	mutex    sync.RWMutex
	taskFunc func(string) error // 实际的长任务函数
}

func NewPreheatingScheduler(taskFunc func(payload string) error) *PreheatingScheduler {
	return &PreheatingScheduler{
		tasks:    make(map[string]*time.Timer),
		taskFunc: taskFunc,
	}
}

func (s *PreheatingScheduler) ScheduleTask(path string, delay time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果已有定时器，先取消
	if timer, exists := s.tasks[path]; exists {
		timer.Stop()
	}

	// 创建新的定时器
	timer := time.AfterFunc(delay, func() {
		// 执行预热任务
		go func() {
			if err := s.taskFunc(path); err != nil {
				log.Printf("预热任务失败: %v", err)
			}

			// 清理定时器
			s.mutex.Lock()
			delete(s.tasks, path)
			s.mutex.Unlock()
		}()
	})

	s.tasks[path] = timer
}
