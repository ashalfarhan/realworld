package service_test

import (
	"log"
	"os"
	"testing"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/db/repository"
	repoMocks "github.com/ashalfarhan/realworld/db/repository/mocks"
	. "github.com/ashalfarhan/realworld/service"
)

var (
	userRepoMock    *repoMocks.UserRepoMock
	articleRepoMock *repoMocks.ArticleRepoMock
	followRepoMock  *repoMocks.FollowingRepoMock
	repo            *repository.Repository
	userService     *UserService
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

	userRepoMock = new(repoMocks.UserRepoMock)
	articleRepoMock = new(repoMocks.ArticleRepoMock)
	followRepoMock = new(repoMocks.FollowingRepoMock)
	repo = &repository.Repository{
		UserRepo:    userRepoMock,
		ArticleRepo: articleRepoMock,
		FollowRepo:  followRepoMock,
	}

	userService = NewUserService(repo)
}

func teardown() {
	log.Println("Tearing down service test")
}
