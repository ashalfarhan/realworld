package service_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestMap map[string]func(*testing.T)

var tests = TestMap{
	"Follow user should fail if already follow": func(t *testing.T) {
		as := assert.New(t)
		userRepoMock.On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()
		followRepoMock.
			On("InsertOne", mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New(repository.ErrDuplicateFollowing)).
			Once()
		u, err := userService.FollowUser(tctx, "username", "username2")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusBadRequest)
		as.Equal(err.Err, ErrAlreadyFollow)
	},
	"Follow user should fail if self follow": func(t *testing.T) {
		username := "username"
		as := assert.New(t)
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{Username: username}, nil).
			Once()
		u, err := userService.FollowUser(tctx, username, "username2")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertNumberOfCalls(t, "InsertOne", 0)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusBadRequest)
		as.Equal(err.Err, ErrSelfFollow)
	},
	"Follow user should fail if not found": func(t *testing.T) {
		as := assert.New(t)

		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(&model.User{}, sql.ErrNoRows).
			Once()
		u, err := userService.FollowUser(tctx, "username", "username2")
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertNumberOfCalls(t, "InsertOne", 0)
		followRepoMock.AssertExpectations(t)

		as.Nil(u)
		as.NotNil(err)
		as.Equal(err.Code, http.StatusNotFound)
		as.Equal(err.Err, ErrNoUserFound)
	},
	"Follow user should success": func(t *testing.T) {
		as := assert.New(t)
		username := "username"

		following := &model.User{
			Username: "username1",
		}
		userRepoMock.
			On("FindOne", mock.Anything, mock.Anything).
			Return(following, nil).
			Once()
		followRepoMock.
			On("InsertOne", mock.Anything, username, following.Username).
			Return(nil).
			Once()

		u, err := userService.FollowUser(tctx, username, following.Username)
		userRepoMock.AssertExpectations(t)
		followRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotNil(u)
		as.True(u.Following)
	},
}

func TestFollowUser(t *testing.T) {
	for name, exec := range tests {
		t.Run(name, exec)
		// Teardown after each subtest
		// like `afterEach` in jest
		followRepoMock.Calls = nil
		userRepoMock.Calls = nil
	}
}
