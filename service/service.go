package service

import (
	"github.com/ashalfarhan/realworld/cache/store"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	UserService    *UserService
	AuthService    *AuthService
	ArticleService *ArticleService
}

func InitService(d *sqlx.DB, s *redis.Client) *Service {
	repo := repository.InitRepository(d)
	store := store.NewCacheStore(s)
	userService := NewUserService(repo)
	articleService := NewArticleService(repo, store)
	authService := NewAuthService(userService)
	return &Service{userService, authService, articleService}
}
