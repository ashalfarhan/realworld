package store

import "github.com/ashalfarhan/realworld/cache"

type CacheStore struct {
	ArticleStore ArticleStore
}

func NewCacheStore() *CacheStore {
	return &CacheStore{
		&ArticleStoreImpl{cache.Ca},
	}
}
