package service

import (
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	UserService    *UserService
	AuthService    *AuthService
	ArticleService *ArticleService
}

func InitService(d *sqlx.DB) *Service {
	repo := repository.InitRepository(d)
	userService := NewUserService(repo)
	articleService := NewArticleService(repo)
	authService := NewAuthService(userService)
	return &Service{userService, authService, articleService}
}
