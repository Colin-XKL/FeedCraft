package util

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

// 预热策略默认情况下每隔8h自动预热一次. 如果36h小时内没有用户请求过来,则不继续进行自动预热

const MAX_PREHEATING_COUNT = 6
const MAX_PREHEATING_GRACE_TIME = 36 * time.Hour
const DEFAULT_PREHEATING_INTERVAL = 8 * time.Hour

type PreheatingContext struct {
	taskKey          string
	firstRequestTime time.Time
	lastRequestTime  time.Time
	preheatingCount  int
	timer            *time.Timer
}

// PreheatingScheduler 预热调度器
type PreheatingScheduler struct {
	contexts map[string]*PreheatingContext
	mutex    sync.RWMutex
	taskFunc func(string) error // 实际的长任务函数
}

func NewPreheatingScheduler(taskFunc func(payload string) error) *PreheatingScheduler {
	return &PreheatingScheduler{
		contexts: make(map[string]*PreheatingContext),
		taskFunc: taskFunc,
	}
}

func (s *PreheatingScheduler) ScheduleTask(recipeName string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ctx, exists := s.contexts[recipeName]
	now := time.Now()

	if !exists {
		ctx = &PreheatingContext{
			taskKey:          recipeName,
			firstRequestTime: now,
			lastRequestTime:  now,
			preheatingCount:  0,
		}
		s.contexts[recipeName] = ctx
	} else {
		ctx.lastRequestTime = now
	}

	// 如果存在旧的定时器，先停止
	if ctx.timer != nil {
		ctx.timer.Stop()
	}

	// 检查是否超过时间窗口或预热次数上限
	if now.Sub(ctx.firstRequestTime) > MAX_PREHEATING_GRACE_TIME ||
		ctx.preheatingCount >= MAX_PREHEATING_COUNT {
		delete(s.contexts, recipeName)
		return
	}

	// 设置下一次预热
	ctx.preheatingCount++
	nextPreheatingTime := ctx.firstRequestTime.Add(time.Duration(ctx.preheatingCount) * DEFAULT_PREHEATING_INTERVAL)
	nextPreheatingTime = nextPreheatingTime.Add(time.Duration(rand.Intn(60)) * time.Second) // 添加一个随机的等待时间,避免短时间大量请求集中

	nextRun := nextPreheatingTime.Sub(now)
	logrus.Debugf("next run after %s", nextRun.String())
	//创建新的定时器
	timer := time.AfterFunc(nextRun, func() {
		// 执行预热任务
		go func() {
			logrus.Infof("running preheating task...(this is [#%d] preheating for key [%s])", ctx.preheatingCount, ctx.taskKey)

			s.mutex.Lock()
			_, taskExist := s.contexts[recipeName]
			s.mutex.Unlock()
			if !taskExist {
				return
			}

			err := s.taskFunc(recipeName)
			if err != nil {
				logrus.Errorf("preheating task for recipe [%s] exec failed. err: %v", recipeName, err)
			}
			s.ScheduleTask(recipeName)
		}()
	})
	ctx.timer = timer
}

type PreheatingTaskInfo struct {
	IsActive        bool
	LastRequestTime time.Time
}

func (s *PreheatingScheduler) GetContextInfo(key string) PreheatingTaskInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ctx, ok := s.contexts[key]
	if !ok {
		return PreheatingTaskInfo{
			IsActive: false,
		}
	}
	return PreheatingTaskInfo{
		IsActive:        true,
		LastRequestTime: ctx.lastRequestTime,
	}
}
