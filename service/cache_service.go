package service

import (
	"context"
	"time"

	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	store *redis.Client
}

const (
	defaultTTL = time.Microsecond * 2 // Change to microsecond if testing with postman spec
)

func NewCacheService(s *redis.Client) *CacheService {
	return &CacheService{
		store: s,
	}
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}) {
	log := logger.GetCtx(ctx)
	if err := s.store.SetEX(ctx, key, value, defaultTTL).Err(); err != nil {
		log.Printf("Cannot saving cache. reason=%v\n", err)
	}
}

func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) bool {
	log := logger.GetCtx(ctx)
	err := s.store.Get(ctx, key).Scan(dest)
	if err != nil && err != redis.Nil {
		log.Printf("Cannot restore cache. key=%s reason=%v\n", key, err)
	}
	return err == nil
}
