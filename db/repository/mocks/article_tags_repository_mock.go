package repository_mocks

import (
	"context"

	. "github.com/ashalfarhan/realworld/db/repository"
	"github.com/stretchr/testify/mock"
)

type ArticleTagsRepoMock struct {
	mock.Mock
}

func (m *ArticleTagsRepoMock) InsertBulk(ctx context.Context, tags []InsertArticleTagsArgs) error {
	args := m.Called(ctx, tags)
	return args.Error(0)
}

func (m *ArticleTagsRepoMock) FindArticleTagsByID(ctx context.Context, articleID string) ([]string, error) {
	args := m.Called(ctx, articleID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *ArticleTagsRepoMock) FindAllTags(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}
