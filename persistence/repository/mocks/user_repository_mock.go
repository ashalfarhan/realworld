package repository_mocks

import (
	"context"

	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/model"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) InsertOne(ctx context.Context, u *model.RegisterUserFields) (*model.User, error) {
	arg := m.Called(ctx, u)
	return arg.Get(0).(*model.User), arg.Error(1)
}

func (m *UserRepoMock) FindOneByID(ctx context.Context, s string) (*model.User, error) {
	arg := m.Called(ctx, s)
	return arg.Get(0).(*model.User), arg.Error(1)
}

func (m *UserRepoMock) FindOne(ctx context.Context, u *repository.FindOneUserFilter) (*model.User, error) {
	arg := m.Called(ctx, u)
	return arg.Get(0).(*model.User), arg.Error(1)
}

func (m *UserRepoMock) UpdateOne(ctx context.Context, d *model.UpdateUserFields, u *model.User) error {
	arg := m.Called(ctx, d, u)
	return arg.Error(0)
}
