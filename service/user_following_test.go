package service_test

import (
	"context"
	"testing"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/stretchr/testify/mock"
)

func TestFollowUser(t *testing.T) {
	t.Run("Follow user should success", func(t *testing.T) {
		id := "id1"
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()

		followRepoMock.
			On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		userService.FollowUser(context.TODO(), id, "username")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertExpectations(t)
	})
}
