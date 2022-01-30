package service

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/db/repository"
)

type Service struct {
	*UserService
	*AuthService
	*ArticleService
}

func InitService(d *sql.DB) *Service {
	repo := repository.InitRepository(d)

	return &Service{
		NewUserService(repo),
		NewAuthService(),
		NewArticleService(repo),
	}
}
