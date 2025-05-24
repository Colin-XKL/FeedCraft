package controller

import (
	"github.com/sirupsen/logrus"
	"math/rand/v2"
	"sync"
	"time"
)

// PreheatingScheduler 预热调度器
type PreheatingScheduler struct {
	tasks    map[string]*time.Timer // key: path, value: timer
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
	timer := time.AfterFunc(delay+time.Duration(rand.IntN(60))*time.Second, func() {
		// 执行预热任务
		go func() {
			logrus.Info("running preheating task...")
			if err := s.taskFunc(path); err != nil {
				logrus.Errorf("预热任务失败: %v", err)
			}

			// 清理定时器
			s.mutex.Lock()
			delete(s.tasks, path)
			s.mutex.Unlock()
		}()
	})

	s.tasks[path] = timer
}
