package cache

import (
	"context"
	"time"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/go-redis/redis/v8"
)

var Ca *redis.Client

const (
	DefaultTTL = time.Microsecond * 2 // Change to microsecond if testing with postman spec
)

func Init() {
	Ca = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.Co.RedisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := Ca.Ping(ctx).Result(); err != nil {
		defer Ca.Close()
		logger.Log.Panicf("Cannot Ping Redis, Reason: %v", err)
		return
	}

	logger.Log.Println("Successfully initialize redis cache")
}
