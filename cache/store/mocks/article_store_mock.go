package mocks

import (
	"context"

	"github.com/ashalfarhan/realworld/model"
	"github.com/stretchr/testify/mock"
)

type ArticleStoreMock struct {
	mock.Mock
}

func (m *ArticleStoreMock) FindOneBySlug(ctx context.Context, arg1 string, arg2 string) *model.Article {
	args := m.Called(ctx, arg1, arg2)
	return args.Get(0).(*model.Article)
}

func (m *ArticleStoreMock) SaveBySlug(ctx context.Context, arg1 string, arg2 string, arg3 *model.Article) {
	m.Called(ctx, arg1, arg2, arg3)
}
