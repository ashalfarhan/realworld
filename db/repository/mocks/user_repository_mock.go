package repository_mocks

import (
	"context"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) InsertOne(ctx context.Context, u *dto.RegisterUserFields) (*model.User, error) {
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

func (m *UserRepoMock) UpdateOne(ctx context.Context, d *dto.UpdateUserFields, u *model.User) error {
	arg := m.Called(ctx, d, u)
	return arg.Error(0)
}
