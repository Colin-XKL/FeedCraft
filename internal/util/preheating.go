package util

import (
	"context"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// 预热策略默认情况下每隔8h自动预热一次. 如果36h小时内没有用户请求过来,则不继续进行自动预热

const MAX_PREHEATING_COUNT = 6
const MAX_PREHEATING_GRACE_TIME = 36 * time.Hour
const DEFAULT_PREHEATING_INTERVAL = 8 * time.Hour
const DEFAULT_PREHEATING_MAX_CONCURRENCY = 2
const DEFAULT_PREHEATING_QUEUE_SIZE = 1000
const DEFAULT_PREHEATING_JITTER = 60 * time.Second
const DEFAULT_PREHEATING_TASK_TIMEOUT = 10 * time.Minute

type PreheatingContext struct {
	taskKey          string
	firstRequestTime time.Time
	lastRequestTime  time.Time
	preheatingCount  int
	timer            *time.Timer
	timerVersion     uint64
	queued           bool
	running          bool
}

type preheatingTask struct {
	recipeName   string
	timerVersion uint64
}

// PreheatingScheduler 预热调度器
type PreheatingScheduler struct {
	contexts       map[string]*PreheatingContext
	mutex          sync.RWMutex
	taskFunc       func(context.Context, string) error // 实际的长任务函数
	shouldRun      func(string) bool
	onSkip         func(string)
	taskQueue      chan preheatingTask
	timerVersion   atomic.Uint64
	maxConcurrency int
	maxQueueSize   int
	maxCount       int
	graceTime      time.Duration
	interval       time.Duration
	jitter         time.Duration
	taskTimeout    time.Duration
}

type PreheatingSchedulerOption func(*PreheatingScheduler)

func WithPreheatingMaxConcurrency(maxConcurrency int) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if maxConcurrency > 0 {
			s.maxConcurrency = maxConcurrency
		}
	}
}

func WithPreheatingQueueSize(queueSize int) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if queueSize > 0 {
			s.maxQueueSize = queueSize
		}
	}
}

func WithPreheatingInterval(interval time.Duration) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if interval > 0 {
			s.interval = interval
		}
	}
}

func WithPreheatingJitter(jitter time.Duration) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if jitter >= 0 {
			s.jitter = jitter
		}
	}
}

func WithPreheatingMaxCount(maxCount int) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if maxCount > 0 {
			s.maxCount = maxCount
		}
	}
}

func WithPreheatingGraceTime(graceTime time.Duration) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if graceTime > 0 {
			s.graceTime = graceTime
		}
	}
}

func WithPreheatingTaskTimeout(taskTimeout time.Duration) PreheatingSchedulerOption {
	return func(s *PreheatingScheduler) {
		if taskTimeout > 0 {
			s.taskTimeout = taskTimeout
		}
	}
}

func NewPreheatingScheduler(taskFunc func(context.Context, string) error, shouldRun func(string) bool, onSkip func(string), options ...PreheatingSchedulerOption) *PreheatingScheduler {
	s := &PreheatingScheduler{
		contexts:       make(map[string]*PreheatingContext),
		taskFunc:       taskFunc,
		shouldRun:      shouldRun,
		onSkip:         onSkip,
		maxConcurrency: DEFAULT_PREHEATING_MAX_CONCURRENCY,
		maxQueueSize:   DEFAULT_PREHEATING_QUEUE_SIZE,
		maxCount:       MAX_PREHEATING_COUNT,
		graceTime:      MAX_PREHEATING_GRACE_TIME,
		interval:       DEFAULT_PREHEATING_INTERVAL,
		jitter:         DEFAULT_PREHEATING_JITTER,
		taskTimeout:    DEFAULT_PREHEATING_TASK_TIMEOUT,
	}

	envClient := GetEnvClient()
	if envClient != nil {
		if envLimit := envClient.GetInt("PREHEATING_MAX_CONCURRENCY"); envLimit > 0 {
			s.maxConcurrency = envLimit
		}
		if envQueueSize := envClient.GetInt("PREHEATING_QUEUE_SIZE"); envQueueSize > 0 {
			s.maxQueueSize = envQueueSize
		}
		if envTimeout := envClient.GetString("PREHEATING_TASK_TIMEOUT"); envTimeout != "" {
			timeout, err := time.ParseDuration(envTimeout)
			if err != nil || timeout <= 0 {
				logrus.Warnf("Invalid FC_PREHEATING_TASK_TIMEOUT %q; using default %s", envTimeout, s.taskTimeout)
			} else {
				s.taskTimeout = timeout
			}
		}
	}
	for _, option := range options {
		option(s)
	}
	s.taskQueue = make(chan preheatingTask, s.maxQueueSize)
	for i := 0; i < s.maxConcurrency; i++ {
		go s.worker()
	}
	logrus.Infof("Preheating scheduler initialized with max concurrency: %d, queue size: %d", s.maxConcurrency, s.maxQueueSize)
	return s
}

func (s *PreheatingScheduler) ScheduleTask(recipeName string) {
	s.scheduleTask(recipeName, true)
}

func (s *PreheatingScheduler) scheduleTask(recipeName string, fromUserRequest bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ctx, exists := s.contexts[recipeName]
	now := time.Now()

	if !exists {
		if !fromUserRequest {
			return
		}
		ctx = &PreheatingContext{
			taskKey:          recipeName,
			firstRequestTime: now,
			lastRequestTime:  now,
			preheatingCount:  0,
		}
		s.contexts[recipeName] = ctx
	} else if fromUserRequest {
		ctx.lastRequestTime = now
		ctx.preheatingCount = 0
		ctx.queued = false
	}

	// 如果存在旧的定时器，先停止
	if ctx.timer != nil {
		ctx.timer.Stop()
		ctx.timer = nil
	}

	// 检查是否超过时间窗口或预热次数上限
	if now.Sub(ctx.lastRequestTime) > s.graceTime ||
		ctx.preheatingCount >= s.maxCount {
		delete(s.contexts, recipeName)
		return
	}

	// 设置下一次预热
	nextRun := s.interval + s.randomJitter()
	timerVersion := s.timerVersion.Add(1)
	ctx.timerVersion = timerVersion
	logrus.Debugf("next run after %s", nextRun.String())
	//创建新的定时器
	timer := time.AfterFunc(nextRun, func() {
		s.enqueueTask(recipeName, timerVersion)
	})
	ctx.timer = timer
}

func (s *PreheatingScheduler) randomJitter() time.Duration {
	if s.jitter <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(s.jitter)))
}

func (s *PreheatingScheduler) enqueueTask(recipeName string, timerVersion uint64) {
	s.mutex.Lock()
	ctx, taskExist := s.contexts[recipeName]
	if !taskExist || ctx.timerVersion != timerVersion {
		s.mutex.Unlock()
		return
	}
	ctx.timer = nil
	if ctx.running || ctx.queued {
		s.mutex.Unlock()
		logrus.Debugf("skip enqueueing duplicate preheating task for recipe [%s]", recipeName)
		return
	}
	ctx.queued = true
	s.mutex.Unlock()

	s.taskQueue <- preheatingTask{
		recipeName:   recipeName,
		timerVersion: timerVersion,
	}
}

func (s *PreheatingScheduler) worker() {
	for task := range s.taskQueue {
		if !s.markTaskRunning(task) {
			continue
		}
		s.runTask(task.recipeName)
	}
}

func (s *PreheatingScheduler) markTaskRunning(task preheatingTask) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ctx, taskExist := s.contexts[task.recipeName]
	if !taskExist || ctx.timerVersion != task.timerVersion || !ctx.queued {
		return false
	}
	if ctx.running {
		ctx.queued = false
		return false
	}
	ctx.queued = false
	ctx.running = true
	ctx.preheatingCount++
	logrus.Infof("running preheating task...(this is [#%d] preheating for key [%s])", ctx.preheatingCount, ctx.taskKey)
	return true
}

func (s *PreheatingScheduler) runTask(recipeName string) {
	shouldReschedule := true
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("preheating task for recipe [%s] panicked: %v", recipeName, r)
		}
		s.finishTask(recipeName, shouldReschedule)
	}()

	if s.shouldRun != nil && !s.shouldRun(recipeName) {
		if s.onSkip != nil {
			s.onSkip(recipeName)
		}
		shouldReschedule = false
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.taskTimeout)
	defer cancel()
	err := s.taskFunc(ctx, recipeName)
	if err != nil {
		logrus.Errorf("preheating task for recipe [%s] exec failed. err: %v", recipeName, err)
	}
}

func (s *PreheatingScheduler) finishTask(recipeName string, shouldReschedule bool) {
	s.mutex.Lock()
	ctx, taskExist := s.contexts[recipeName]
	if !taskExist {
		s.mutex.Unlock()
		return
	}
	ctx.running = false
	if !shouldReschedule {
		delete(s.contexts, recipeName)
		s.mutex.Unlock()
		return
	}
	s.mutex.Unlock()

	s.scheduleTask(recipeName, false)
}

type PreheatingTaskInfo struct {
	IsActive        bool
	LastRequestTime time.Time
}

func (s *PreheatingScheduler) GetContextInfo(key string) PreheatingTaskInfo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
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
