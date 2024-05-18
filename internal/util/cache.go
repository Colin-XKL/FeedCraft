package util

import (
	"context"
	"github.com/redis/go-redis/v9"
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
