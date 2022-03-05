package service_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/repository"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	userService := NewUserService(repo)
	testCases := []struct {
		desc       string
		mockReturn error
		errCode    int
		errError   error
	}{
		{
			desc:       "Register should fail if username exist",
			mockReturn: errors.New(repository.ErrDuplicateUsername),
			errCode:    http.StatusBadRequest,
			errError:   ErrUsernameExist,
		},
		{
			desc:       "Register should fail if email exist",
			mockReturn: errors.New(repository.ErrDuplicateEmail),
			errCode:    http.StatusBadRequest,
			errError:   ErrEmailExist,
		},
		{
			desc:       "Register should fail if db error",
			mockReturn: sql.ErrTxDone,
			errCode:    http.StatusInternalServerError,
			errError:   conduit.ErrInternal,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			as := assert.New(t)
			userRepoMock.On("InsertOne", mock.Anything, mock.Anything).Return(tC.mockReturn).Once()
			_, err := userService.Register(context.TODO(), &RegisterArgs{})
			userRepoMock.AssertExpectations(t)

			as.NotNil(err)
			as.Equal(err.Code, tC.errCode)
			as.Equal(err.Error, tC.errError)
		})
	}

	t.Run("Register should success", func(t *testing.T) {
		pw := "password"
		as := assert.New(t)
		userRepoMock.On("InsertOne", mock.Anything, mock.Anything).Return(nil).Once()
		reg, err := userService.Register(context.TODO(), &RegisterArgs{
			Password: pw,
		})
		userRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotEqual(pw, reg.Password, "Registered user password should be hashed")
	})
}
