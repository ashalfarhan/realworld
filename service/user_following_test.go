package service_test

import (
	"context"
	"errors"
	"database/sql"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollowUser(t *testing.T) {
	t.Run("Follow user should fail if already follow", func(t *testing.T) {
		t.Cleanup(func() {
			followRepoMock.Calls = nil
			userRepoMock.Calls = nil
		})
		as := assert.New(t)
		userRepoMock.On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()
		followRepoMock.
			On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New(repository.ErrDuplicateFollowing)).
			Once()
		u, err := userService.FollowUser(context.TODO(), "uid", "username")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusBadRequest)
		as.Equal(err.Error, ErrAlreadyFollow)
	})

	t.Run("Follow user should fail if self follow", func(t *testing.T) {
		t.Cleanup(func() {
			followRepoMock.Calls = nil
			userRepoMock.Calls = nil
		})
		uid := "uid"
		as := assert.New(t)
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{ID: uid}, nil).
			Once()
		u, err := userService.FollowUser(context.TODO(), uid, "username")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertNotCalled(t, "InsertOne", mock.Anything, mock.Anything, mock.Anything)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusBadRequest)
		as.Equal(err.Error, ErrSelfFollow)
	})

	t.Run("Follow user should fail if not found", func(t *testing.T) {
		t.Cleanup(func() {
			followRepoMock.Calls = nil
			userRepoMock.Calls = nil
		})
		as := assert.New(t)
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, sql.ErrNoRows).
			Once()
		u, err := userService.FollowUser(context.TODO(), "id", "username")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertNotCalled(t, "InsertOne", mock.Anything, mock.Anything, mock.Anything)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusNotFound)
		as.Equal(err.Error, ErrNoUserFound)
	})

	t.Run("Follow user should success", func(t *testing.T) {
		t.Cleanup(func() {
			followRepoMock.Calls = nil
			userRepoMock.Calls = nil
		})
		as := assert.New(t)
		id := "id1"
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()
		followRepoMock.
			On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		u, err := userService.FollowUser(context.TODO(), id, "username")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotNil(u)
		as.True(u.Following)
	})
}
