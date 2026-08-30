package util

import (
	"FeedCraft/internal/constant"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
)

var cachedFuncFlight singleflight.Group

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
// Concurrent misses for the same key share one in-flight valFunc (singleflight) to avoid cache stampede.
func CachedFuncWithPreLog(cacheKey string, valFunc func() (string, error), preLog func(isCached bool)) (string, error) {
	result, err, _ := cachedFuncFlight.Do(cacheKey, func() (interface{}, error) {
		cached, getErr := CacheGetString(cacheKey)
		isCached := getErr == nil && cached != ""

		if preLog != nil {
			preLog(isCached)
		}

		if isCached {
			return cached, nil
		}

		processedContent, getValErr := valFunc()
		if getValErr != nil {
			return "", getValErr
		}
		if cacheErr := CacheSetString(cacheKey, processedContent, constant.WebContentExpire); cacheErr != nil {
			logrus.Warn("failed to cache result")
		}
		return processedContent, nil
	})
	if err != nil {
		return "", err
	}
	value, _ := result.(string)
	return value, nil
}

// CachedFunc 先尝试取缓存, 如不存在, 则调用valFunc 获取值并写入缓存
func CachedFunc(cacheKey string, valFunc func() (string, error)) (string, error) {
	return CachedFuncWithPreLog(cacheKey, valFunc, nil)
}
