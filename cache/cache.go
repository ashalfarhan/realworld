package cache

import (
	"context"
	"time"

	"github.com/ashalfarhan/realworld/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func Init() *redis.Client {
	Ca := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.RedisPass,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := Ca.Ping(ctx).Result(); err != nil {
		defer Ca.Close()
		logrus.Panicf("Cannot Ping Redis, Reason: %v", err)
		return nil
	}
	logrus.Println("Successfully initialize redis cache")
	return Ca
}
