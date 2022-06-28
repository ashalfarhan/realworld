package service_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	. "github.com/ashalfarhan/realworld/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterFail(t *testing.T) {
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			as := assert.New(t)

			userRepoMock.On("InsertOne", mock.Anything, mock.Anything).Return(&model.User{}, tC.mockReturn).Once()
			d, err := userService.Insert(tctx, &model.RegisterUserFields{})
			userRepoMock.AssertExpectations(t)

			as.Nil(d)
			if as.NotNil(err) {
				as.Equal(err.Code, tC.errCode)
				as.Equal(err.Err, tC.errError)
			}
		})
	}
}

func TestRegisterSuccess(t *testing.T) {
	pw := "password"
	as := assert.New(t)

	userRepoMock.On("InsertOne", mock.Anything, mock.Anything).Return(&model.User{}, nil).Once()
	reg, err := userService.Insert(tctx, &model.RegisterUserFields{Password: pw})
	userRepoMock.AssertExpectations(t)

	as.Nil(err)
	if as.NotNil(reg) {
		as.NotEqual(pw, reg.Password, "Registered user password should be hashed")
	}
}

func TestUpdateFail(t *testing.T) {
	testCases := []struct {
		desc       string
		mockReturn error
		errCode    int
		errError   error
	}{
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			as := assert.New(t)

			userRepoMock.On("FindOneByUsername", mock.Anything, mock.Anything).Return(&model.User{}, nil)
			userRepoMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(tC.mockReturn).Once()
			d, err := userService.Update(tctx, &model.UpdateUserFields{}, "")
			userRepoMock.AssertExpectations(t)

			as.Nil(d)
			if as.NotNil(err) {
				as.Equal(err.Code, tC.errCode)
				as.Equal(err.Err, tC.errError)
			}
		})
	}
}
