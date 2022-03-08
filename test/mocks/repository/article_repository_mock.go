package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/stretchr/testify/mock"
)

type ArticleRepoMock struct {
	mock.Mock
}

func (m *ArticleRepoMock) InsertOne(ctx context.Context, a *dto.CreateArticleFields, slug string) (*model.Article, error) {
	args := m.Called(ctx, a, slug)
	return args.Get(0).(*model.Article), args.Error(1)
}

func (m *ArticleRepoMock) FindOneBySlug(ctx context.Context, s string) (*model.Article, error) {
	args := m.Called(ctx, s)
	return args.Get(1).(*model.Article), args.Error(1)
}

func (m *ArticleRepoMock) DeleteBySlug(ctx context.Context, s string) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *ArticleRepoMock) UpdateOneBySlug(ctx context.Context, d *dto.UpdateArticleFields, a *model.Article) error {
	args := m.Called(ctx, d, a)
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
