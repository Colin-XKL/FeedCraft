package util

import (
	"FeedCraft/internal/constant"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

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

func CacheSetString(key string, value string, ttl time.Duration) error {
	rdb := GetRedisClient()
	return rdb.Set(context.Background(), key, value, ttl).Err()
}
func CacheGetString(key string) (string, error) {
	rdb := GetRedisClient()
	return rdb.Get(context.Background(), key).Result()
}

// CachedFunc 先尝试取缓存, 如不存在, 则调用valFunc 获取值并写入缓存
func CachedFunc(cacheKey string, valFunc func() (string, error)) (string, error) {
	final := ""
	cached, err := CacheGetString(cacheKey)
	if err != nil || cached == "" {
		processedContent, getValErr := valFunc()
		if getValErr != nil {
			return "", getValErr
		} else {
			final = processedContent
			cacheErr := CacheSetString(cacheKey, processedContent, constant.WebContentExpire)
			if cacheErr != nil {
				logrus.Warn("failed to cache result")
				//logrus.Warnf("failed to cache result of craft [%s] for article [%s], %v\n", craftName,
				//	originalTitle, cacheErr)
			}
		}
	} else {
		final = cached
	}

	return final, nil
}
