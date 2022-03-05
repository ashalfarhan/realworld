package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/stretchr/testify/mock"
)

type ArticleRepoMock struct {
	mock.Mock
}

func (m *ArticleRepoMock) InsertOne(ctx context.Context, a *model.Article) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *ArticleRepoMock) FindOneBySlug(ctx context.Context, a *model.Article) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *ArticleRepoMock) DeleteBySlug(ctx context.Context, s string) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *ArticleRepoMock) UpdateOneBySlug(ctx context.Context, s string, uv *repository.UpdateArticleValues, a *model.Article) error {
	args := m.Called(ctx, s, uv, a)
	return args.Error(0)
}

func (m *ArticleRepoMock) Find(ctx context.Context, a *repository.FindArticlesArgs) (model.Articles, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(model.Articles), args.Error(1)
}

func (m *ArticleRepoMock) FindByFollowed(ctx context.Context, a *repository.FindArticlesArgs) (model.Articles, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(model.Articles), args.Error(1)
}
