package service_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/db/repository"
	repoMocks "github.com/ashalfarhan/realworld/db/repository/mocks"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/mock"
)

var (
	userRepoMock        *repoMocks.UserRepoMock
	articleRepoMock     *repoMocks.ArticleRepoMock
	followRepoMock      *repoMocks.FollowingRepoMock
	articleTagsRepoMock *repoMocks.ArticleTagsRepoMock
	repo                *repository.Repository
	userService         *UserService
	articleService      *ArticleService
	tctx                = context.TODO()
	mockCtx             = mock.Anything
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	config.Co.Env = "test"
	log.Println("Setting up service test")

	userRepoMock = new(repoMocks.UserRepoMock)
	articleRepoMock = new(repoMocks.ArticleRepoMock)
	followRepoMock = new(repoMocks.FollowingRepoMock)
	articleTagsRepoMock = new(repoMocks.ArticleTagsRepoMock)
	repo = &repository.Repository{
		UserRepo:        userRepoMock,
		ArticleRepo:     articleRepoMock,
		FollowRepo:      followRepoMock,
		ArticleTagsRepo: articleTagsRepoMock,
	}

	userService = NewUserService(repo)
	articleService = NewArticleService(repo)
}
