package repository_mocks

import (
	"context"

	"github.com/ashalfarhan/realworld/model"
	"github.com/stretchr/testify/mock"
)

type CommentRepoMock struct {
	mock.Mock
}

func (m *CommentRepoMock) InsertOne(ctx context.Context, c *model.Comment) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CommentRepoMock) FindByArticleID(ctx context.Context, articleID string) ([]*model.Comment, error) {
	args := m.Called(ctx, articleID)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *CommentRepoMock) DeleteByID(ctx context.Context, commentID string) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *CommentRepoMock) FindOneByID(ctx context.Context, commentID string) (*model.Comment, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.Comment), args.Error(1)
}
