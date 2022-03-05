package service_test

import (
	"log"
	"os"
	"testing"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/db/repository"
)

var (
	userRepoMock    *repository.UserRepoMock
	articleRepoMock *repository.ArticleRepoMock
	followRepoMock  *repository.FollowingRepoMock
	repo            *repository.Repository
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	config.Co.Env = "test"
	log.Println("Setting up service test")

	userRepoMock = new(repository.UserRepoMock)
	articleRepoMock = new(repository.ArticleRepoMock)
	followRepoMock = new(repository.FollowingRepoMock)
	repo = &repository.Repository{
		UserRepo:    userRepoMock,
		ArticleRepo: articleRepoMock,
		FollowRepo:  followRepoMock,
	}
}

func teardown() {
	log.Println("Tearing down service test")
}
