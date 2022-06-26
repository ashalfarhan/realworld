package store

import (
	"context"
	"fmt"

	"github.com/ashalfarhan/realworld/cache"
	"github.com/ashalfarhan/realworld/model"
	"github.com/go-redis/redis/v8"
)

type ArticleStoreImpl struct {
	client *redis.Client
}

type ArticleStore interface {
	FindOneBySlug(context.Context, string, string) *model.Article
	SaveBySlug(context.Context, string, string, *model.Article)
}

var prefix = "articles"

func (s *ArticleStoreImpl) FindOneBySlug(ctx context.Context, slug, userID string) *model.Article {
	key := fmt.Sprintf("%s|slug:%s|user_id:%s", prefix, slug, userID)
	res := new(model.Article)
	if err := s.client.Get(ctx, key).Scan(res); err != nil {
		return nil
	}
	return res
}

func (s *ArticleStoreImpl) SaveBySlug(ctx context.Context, slug, userID string, a *model.Article) {
	key := fmt.Sprintf("%s|slug:%s|user_id:%s", prefix, slug, userID)
	s.client.SetEX(ctx, key, a, cache.DefaultTTL)
}
