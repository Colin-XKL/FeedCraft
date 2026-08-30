package util

import (
	"FeedCraft/internal/constant"
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var cacheFlight = newCacheFlightGroup()

type cacheFlightCall struct {
	ctx       context.Context
	cancel    context.CancelFunc
	done      chan struct{}
	value     string
	err       error
	waiters   int
	completed bool
}

type cacheFlightGroup struct {
	mu    sync.Mutex
	calls map[string]*cacheFlightCall
}

func newCacheFlightGroup() *cacheFlightGroup {
	return &cacheFlightGroup{calls: make(map[string]*cacheFlightCall)}
}

func (g *cacheFlightGroup) do(
	ctx context.Context,
	key string,
	fn func(context.Context) (string, error),
) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	g.mu.Lock()
	call := g.calls[key]
	if call != nil {
		call.waiters++
	} else {
		sharedCtx, cancel := context.WithCancel(context.Background())
		call = &cacheFlightCall{
			ctx:     sharedCtx,
			cancel:  cancel,
			done:    make(chan struct{}),
			waiters: 1,
		}
		g.calls[key] = call
		go g.run(key, call, fn)
	}
	g.mu.Unlock()

	select {
	case <-call.done:
		return call.value, call.err
	default:
	}

	select {
	case <-call.done:
		return call.value, call.err
	case <-ctx.Done():
		g.leave(key, call)
		return "", ctx.Err()
	}
}

func (g *cacheFlightGroup) run(
	key string,
	call *cacheFlightCall,
	fn func(context.Context) (string, error),
) {
	value, err := fn(call.ctx)

	g.mu.Lock()
	call.value = value
	call.err = err
	call.completed = true
	if g.calls[key] == call {
		delete(g.calls, key)
	}
	close(call.done)
	call.cancel()
	g.mu.Unlock()
}

func (g *cacheFlightGroup) leave(key string, call *cacheFlightCall) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if call.completed || g.calls[key] != call {
		return
	}
	call.waiters--
	if call.waiters == 0 {
		delete(g.calls, key)
		call.cancel()
	}
}

// GetRedisClient 返回一个非空的redis client
func GetRedisClient() *redis.Client {
	envClient := GetEnvClient()
	if envClient == nil {
		log.Fatalf("get env client error.")
		return nil
	}
	redisURI := envClient.GetString("REDIS_URI")

	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		log.Fatalf("parse redis uri fail. err:%v", err)
	}
	rdb := redis.NewClient(opts)
	if rdb == nil {
		log.Fatalf("create redis client error.")
	}
	return rdb
}

// tryGetRedisClient returns a Redis client or an error without fatally crashing.
// It is safe to call in environments where Redis may not be configured (e.g. tests).
func tryGetRedisClient() (*redis.Client, error) {
	envClient := GetEnvClient()
	if envClient == nil {
		return nil, errors.New("env client not available")
	}
	redisURI := strings.TrimSpace(envClient.GetString("REDIS_URI"))
	if redisURI == "" {
		return nil, errors.New("REDIS_URI is not configured")
	}
	opts, err := redis.ParseURL(redisURI)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opts)
	if rdb == nil {
		return nil, errors.New("failed to create redis client")
	}
	return rdb, nil
}

// TryCacheSetString is like CacheSetString but returns an error instead of crashing
// when Redis is not configured. Safe to call in environments without Redis.
func TryCacheSetString(key string, value string, ttl time.Duration) error {
	rdb, err := tryGetRedisClient()
	if err != nil {
		return err
	}
	return rdb.Set(context.Background(), key, value, ttl).Err()
}

// TryCacheGetString is like CacheGetString but returns an error instead of crashing
// when Redis is not configured. Safe to call in environments without Redis.
func TryCacheGetString(key string) (string, error) {
	rdb, err := tryGetRedisClient()
	if err != nil {
		return "", err
	}
	return rdb.Get(context.Background(), key).Result()
}

func CacheSetString(key string, value string, ttl time.Duration) error {
	rdb := GetRedisClient()
	return rdb.Set(context.Background(), key, value, ttl).Err()
}
func CacheGetString(key string) (string, error) {
	rdb := GetRedisClient()
	return rdb.Get(context.Background(), key).Result()
}

// CachedFuncWithPreLog tries to get from cache, invokes preLog if provided, and if absent, calls valFunc and saves to cache.
// Concurrent misses for the same key are coalesced to avoid cache stampede.
func CachedFuncWithPreLog(cacheKey string, valFunc func() (string, error), preLog func(isCached bool)) (string, error) {
	return CachedFuncWithPreLogContext(context.Background(), cacheKey, func(context.Context) (string, error) {
		return valFunc()
	}, preLog)
}

// CachedFuncWithPreLogContext coalesces same-key cache misses while keeping
// each caller's cancellation independent. Shared work stops when all waiters leave.
func CachedFuncWithPreLogContext(
	ctx context.Context,
	cacheKey string,
	valFunc func(context.Context) (string, error),
	preLog func(isCached bool),
) (string, error) {
	return cachedFuncWithStoreContext(ctx, cacheKey, CacheGetString, func(value string) error {
		return CacheSetString(cacheKey, value, constant.WebContentExpire)
	}, valFunc, preLog)
}

func cachedFuncWithStore(
	cacheKey string,
	get func(string) (string, error),
	set func(string) error,
	valFunc func() (string, error),
	preLog func(isCached bool),
) (string, error) {
	return cachedFuncWithStoreContext(context.Background(), cacheKey, get, set, func(context.Context) (string, error) {
		return valFunc()
	}, preLog)
}

func cachedFuncWithStoreContext(
	ctx context.Context,
	cacheKey string,
	get func(string) (string, error),
	set func(string) error,
	valFunc func(context.Context) (string, error),
	preLog func(isCached bool),
) (string, error) {
	cached, err := get(cacheKey)
	isCached := err == nil && cached != ""
	if preLog != nil {
		preLog(isCached)
	}
	if isCached {
		return cached, nil
	}

	return cacheFlight.do(ctx, cacheKey, func(sharedCtx context.Context) (string, error) {
		cached, err := get(cacheKey)
		if err == nil && cached != "" {
			return cached, nil
		}
		processedContent, getValErr := valFunc(sharedCtx)
		if getValErr != nil {
			return "", getValErr
		}
		if cacheErr := set(processedContent); cacheErr != nil {
			logrus.Warn("failed to cache result")
		}
		return processedContent, nil
	})
}

// CachedFunc 先尝试取缓存, 如不存在, 则调用valFunc 获取值并写入缓存
func CachedFunc(cacheKey string, valFunc func() (string, error)) (string, error) {
	return CachedFuncWithPreLog(cacheKey, valFunc, nil)
}

// CachedFuncContext is CachedFunc with independent per-caller cancellation.
func CachedFuncContext(ctx context.Context, cacheKey string, valFunc func(context.Context) (string, error)) (string, error) {
	return CachedFuncWithPreLogContext(ctx, cacheKey, valFunc, nil)
}
