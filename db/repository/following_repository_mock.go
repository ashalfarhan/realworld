package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FollowingRepoMock struct {
	mock.Mock
}

func (m *FollowingRepoMock) InsertOne(ctx context.Context, s string, sa string) error {
	args := m.Called(ctx, s, sa)
	return args.Error(0)
}

func (m *FollowingRepoMock) DeleteOneIDs(ctx context.Context, s string, sa string) error {
	args := m.Called(ctx, s, sa)
	return args.Error(0)
}

func (m *FollowingRepoMock) FindOneByIDs(ctx context.Context, s string, sa string) (*string, error) {
	args := m.Called(ctx, s, sa)
	return args.Get(0).(*string), args.Error(1)
}
