package cache

import (
	"context"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/config"
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
		conduit.Logger.Panicf("Cannot Ping Redis, Reason: %v", err)
		return
	}

	conduit.Logger.Println("Successfully initialize redis cache")
}
