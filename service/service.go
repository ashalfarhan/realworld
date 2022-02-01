package service

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/db/repository"
)

type Service struct {
	UserService    *UserService
	AuthService    *AuthService
	ArticleService *ArticleService
}

func InitService(d *sql.DB) *Service {
	repo := repository.InitRepository(d)

	return &Service{
		NewUserService(repo),
		NewAuthService(),
		NewArticleService(repo),
	}
}
