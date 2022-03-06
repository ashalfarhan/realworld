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
	us := NewUserService(repo)
	as := NewArticleService(repo)
	return &Service{
		us,
		NewAuthService(us),
		as,
	}
}
