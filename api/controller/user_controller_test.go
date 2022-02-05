package controller

import (
	"net/http"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/test"
)

type DtoError map[string]map[string][]string

var userController = &UserController{}

func TestLoginController(t *testing.T) {
	t.Run("Should response error if not provide email or username", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "secret",
			},
		}
		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/login", payload, userController.LoginUser, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["username"]) == 0 && len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in username and email, but got %v", body["errors"])
		}
	})
}

func TestRegisterController(t *testing.T) {
	t.Run("Should response error if not provide email and username", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "secret",
			},
		}

		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, userController.RegisterUser, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["username"]) == 0 && len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in username and email, but got %v", body["errors"])
		}
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
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, userController.RegisterUser, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if body["errors"]["password"][0] != "min 8" {
			t.Fatalf("expected validation errors in password, but got %v", body["errors"])
		}
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
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, userController.RegisterUser, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in email, but got %v", body["errors"])
		}
	})
}

func TestUpdateUserController(t *testing.T) {
	t.Run("Should response error if new password length is not the min (8)", func(t *testing.T) {
		payload := response.M{
			"user": response.M{
				"password": "asd",
			},
		}
		var body DtoError
		res := test.MakeRequest(http.MethodPost, "/api/users/register", payload, userController.RegisterUser, &body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %d", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if body["errors"]["password"][0] == "min 8 " {
			t.Fatalf("expected validation errors in password, but got %v", body["errors"])
		}
	})
}
