package cache

import (
	"context"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/go-redis/redis/v8"
)

var Ca *redis.Client

func Init() {
	Ca = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.Co.RedisPass,
		DB:       0,
	})

	if _, err := Ca.Ping(context.TODO()).Result(); err != nil {
		defer Ca.Close()
		logger.Log.Panicf("Cannot Ping Redis, Reason: %v", err)
		return
	}

	logger.Log.Println("Successfully initialize redis cache")
}
