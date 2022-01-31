package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
)

func TestLoginController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if not provide email or username", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"password": "secret",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.LoginUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})
}

func TestRegisterController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if not provide email and username", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"password": "secret",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})

	t.Run("Should response error if not password is less than 8", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"email":    "johndoe@mail.com",
			"username": "doejohn",
			"password": "secret",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})

	t.Run("Should response error if not a valid email", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"email":    "johndoe@mail@~asdasd.com",
			"username": "doejohn",
			"password": "secret2022",
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})
}

func TestUpdateUserController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if new password length is not the min (8)", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"password": "asd",
		})

		req := httptest.NewRequest(http.MethodPut, "/api/user", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.UpdateCurrentUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

	})
}
