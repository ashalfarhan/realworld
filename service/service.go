package service

import (
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	UserService    *UserService
	AuthService    *AuthService
	ArticleService *ArticleService
}

func InitService(d *sqlx.DB) *Service {
	repo := repository.InitRepository(d)

	return &Service{
		NewUserService(repo),
		NewAuthService(),
		NewArticleService(repo),
	}
}
