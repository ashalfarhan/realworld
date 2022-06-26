package store

import "github.com/go-redis/redis/v8"

type CacheStore struct {
	ArticleStore ArticleStore
}

func NewCacheStore(c *redis.Client) *CacheStore {
	return &CacheStore{
		&ArticleStoreImpl{c},
	}
}
