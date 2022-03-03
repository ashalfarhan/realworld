package controller

import (
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/test"
	"github.com/stretchr/testify/assert"
)

type DtoError map[string]map[string][]string

var userController = &UserController{}

func TestLoginController(t *testing.T) {
	t.Parallel()

	t.Run("Should response error if not provide email or username", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "secret",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/login", payload, &body, userController.LoginUser)
		as := assert.New(t)
		as.Equal(res.StatusCode, http.StatusUnprocessableEntity)
		as.Greater(len(body["errors"]["username"]), 0, "validation error in username")
		as.Greater(len(body["errors"]["email"]), 0, "validation error in email")
		as.Contains(body["errors"]["username"][0], "required")
		as.Contains(body["errors"]["email"][0], "required")
	})
}

func TestRegisterController(t *testing.T) {
	t.Parallel()

	t.Run("Should response error if not provide email and username", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "secret2020",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, &body, userController.RegisterUser)
		as := assert.New(t)
		as.Equal(res.StatusCode, http.StatusUnprocessableEntity)
		as.Greater(len(body["errors"]["username"]), 0, "validation error in username")
		as.Greater(len(body["errors"]["email"]), 0, "validation error in email")
		as.Contains(body["errors"]["username"][0], "required")
		as.Contains(body["errors"]["email"][0], "required")
	})

	t.Run("Should response error if password is less than 8", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"email":    "johndoe@mail.com",
				"username": "doejohn",
				"password": "secret",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, &body, userController.RegisterUser)
		as := assert.New(t)
		as.Equal(res.StatusCode, http.StatusUnprocessableEntity)
		as.Greater(len(body["errors"]["password"]), 0, "validation error in password")
		as.Contains(body["errors"]["password"][0], "min 8")
	})

	t.Run("Should response error if not a valid email", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"email":    "johndoe@mail@~asdasd.com",
				"username": "doejohn",
				"password": "secret2022",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, &body, userController.RegisterUser)
		as := assert.New(t)
		as.Equal(res.StatusCode, http.StatusUnprocessableEntity)
		as.Greater(len(body["errors"]["email"]), 0, "validation error in email")
		as.Contains(body["errors"]["email"][0], "email")
	})
}

func TestUpdateUserController(t *testing.T) {
	t.Parallel()

	t.Run("Should response error if new password length is not the min (8)", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "asd",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, &body, userController.RegisterUser)
		as := assert.New(t)
		as.Equal(res.StatusCode, http.StatusUnprocessableEntity)
		as.Greater(len(body["errors"]["password"]), 0, "validation error in password")
		as.Contains(body["errors"]["password"][0], "min 8")
	})
}
