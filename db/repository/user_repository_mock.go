package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) InsertOne(ctx context.Context, u *model.User) error {
	arg := m.Called(ctx, u)
	return arg.Error(0)
}

func (m *UserRepoMock) FindOneByID(ctx context.Context, s string, u *model.User) error {
	arg := m.Called(ctx, s, u)
	return arg.Error(0)
}

func (m *UserRepoMock) FindOne(ctx context.Context, u *model.User) error {
	arg := m.Called(ctx, u)
	return arg.Error(0)
}

func (m *UserRepoMock) UpdateOne(ctx context.Context, uv *UpdateUserValues) error {
	arg := m.Called(ctx, uv)
	return arg.Error(0)
}
