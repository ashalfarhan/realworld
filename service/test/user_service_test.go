package service_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase struct {
	desc       string
	mockReturn error
	errCode    int
	errError   error
}

func TestRegister(t *testing.T) {
	testCases := []TestCase{
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
			userRepoMock.
				On("InsertOne", mock.Anything, mock.Anything).
				Return(&model.User{}, tC.mockReturn).
				Once()
			_, err := userService.Insert(tctx, &dto.RegisterUserFields{})
			userRepoMock.AssertExpectations(t)

			as.NotNil(err)
			as.Equal(err.Code, tC.errCode)
			as.Equal(err.Error, tC.errError)
		})
	}

	t.Run("Register should success", func(t *testing.T) {
		pw := "password"
		as := assert.New(t)
		userRepoMock.
			On("InsertOne", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()
		reg, err := userService.Insert(tctx, &dto.RegisterUserFields{
			Password: pw,
		})
		userRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotNil(reg)
		as.NotEqual(pw, reg.Password, "Registered user password should be hashed")
	})
}

func TestGetOneById(t *testing.T) {
	testCases := []TestCase{
		{
			desc:       "Get one by id should fail if no rows",
			mockReturn: sql.ErrNoRows,
			errCode:    http.StatusNotFound,
			errError:   ErrNoUserFound,
		},
		{
			desc:       "Get one by id should fail if db error",
			mockReturn: sql.ErrTxDone,
			errCode:    http.StatusInternalServerError,
			errError:   conduit.ErrInternal,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			as := assert.New(t)
			userRepoMock.
				On("FindOneByID", mock.Anything, mock.Anything).
				Return(&model.User{}, tC.mockReturn).
				Once()
			_, err := userService.GetOneById(tctx, "id")
			userRepoMock.AssertExpectations(t)

			as.NotNil(err)
			as.Equal(err.Code, tC.errCode)
			as.Equal(err.Error, tC.errError)
		})
	}

	t.Run("Get one by id should success", func(t *testing.T) {
		as := assert.New(t)

		userRepoMock.
			On("FindOneByID", mock.Anything, mock.Anything).
			Return(&model.User{}, nil).
			Once()
		u, err := userService.GetOneById(tctx, "id")
		userRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotNil(u)
	})
}

func TestUpdate(t *testing.T) {
	userRepoMock.
		On("FindOneByID", mock.Anything, mock.Anything).
		Return(&model.User{}, nil)
	password := "asd"
	email := "asd@mail.com"
	data := &dto.UpdateUserFields{
		Password: &password,
		Email:    &email,
	}

	testCases := []TestCase{
		{
			desc:       "Update should fail if username exist",
			mockReturn: errors.New(repository.ErrDuplicateUsername),
			errCode:    http.StatusBadRequest,
			errError:   ErrUsernameExist,
		},
		{
			desc:       "Update should fail if email exist",
			mockReturn: errors.New(repository.ErrDuplicateEmail),
			errCode:    http.StatusBadRequest,
			errError:   ErrEmailExist,
		},
		{
			desc:       "Update should fail if db error",
			mockReturn: sql.ErrTxDone,
			errCode:    http.StatusInternalServerError,
			errError:   conduit.ErrInternal,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			as := assert.New(t)

			userRepoMock.
				On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
				Return(tC.mockReturn).
				Once()
			_, err := userService.Update(tctx, data, "")
			userRepoMock.AssertExpectations(t)

			as.NotNil(err)
			as.Equal(err.Code, tC.errCode)
			as.Equal(err.Error, tC.errError)
		})
	}

	t.Run("Update user should success", func(t *testing.T) {
		as := assert.New(t)

		userRepoMock.
			On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		u, err := userService.Update(tctx, data, "")
		userRepoMock.AssertExpectations(t)

		as.Nil(err)
		as.NotEqual(u.Password, password, "Should hash new password")
	})
}
