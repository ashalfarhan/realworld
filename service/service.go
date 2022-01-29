package service

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/db/repository"
)

type Service struct {
	US *UserService
	AS *AuthService
}

func InitService(d *sql.DB) *Service {
	repo := repository.InitRepository(d)

	return &Service{
		US: NewUserService(repo),
		AS: NewAuthService(),
	}
}
