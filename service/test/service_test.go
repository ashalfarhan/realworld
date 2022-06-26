package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/ashalfarhan/realworld/cache/store"
	storeMocks "github.com/ashalfarhan/realworld/cache/store/mocks"
	"github.com/ashalfarhan/realworld/config"
	"github.com/ashalfarhan/realworld/persistence/repository"
	repoMocks "github.com/ashalfarhan/realworld/persistence/repository/mocks"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/stretchr/testify/mock"
)

var (
	userRepoMock        *repoMocks.UserRepoMock
	articleRepoMock     *repoMocks.ArticleRepoMock
	followRepoMock      *repoMocks.FollowingRepoMock
	articleTagsRepoMock *repoMocks.ArticleTagsRepoMock
	repo                *repository.Repository

	articleStoreMock *storeMocks.ArticleStoreMock
	cacheStore       *store.CacheStore

	userService    *UserService
	articleService *ArticleService

	tctx    = context.TODO()
	mockCtx = mock.Anything
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	config.Env = "test"
	logger.Init()

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

	articleStoreMock = new(storeMocks.ArticleStoreMock)
	cacheStore = &store.CacheStore{
		ArticleStore: articleStoreMock,
	}

	userService = NewUserService(repo)
	articleService = NewArticleService(repo, cacheStore)
}
