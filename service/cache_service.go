package service

import (
	"context"
	"time"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type CacheService struct {
	store  *redis.Client
	logger *logrus.Entry
}

const (
	defaultTTL = time.Hour * 2
)

func NewCacheService(s *redis.Client) *CacheService {
	return &CacheService{
		s,
		conduit.NewLogger("Service", "CacheService"),
	}
}

func (s *CacheService) Set(ctx context.Context, key string, value interface{}) {
	if err := s.store.SetEX(ctx, key, value, defaultTTL).Err(); err != nil {
		s.logger.Printf("Cannot saving cache. reason=%v\n", err)
	}
}

func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) bool {
	err := s.store.Get(ctx, key).Scan(dest)
	if err != nil {
		s.logger.Printf("Cannot restore cache. key=%s reason=%v\n", key, err)
	}

	return err == nil
}
