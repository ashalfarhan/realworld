package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if not provide email or username", func(t *testing.T) {
		payload := strings.NewReader(`{"password": "secret"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", payload)
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
		payload := strings.NewReader(`{"password": "secret"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", payload)
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})

	t.Run("Should response error if not password is less than 8", func(t *testing.T) {
		payload := strings.NewReader(`{"email": "johndoe@mail.com", "username": "doejohn", "password": "secret"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", payload)
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})

	t.Run("Should response error if not a valid email", func(t *testing.T) {
		payload := strings.NewReader(`{"email": "johndoe@mail@~asdasd.com", "username": "doejohn", "password": "secret2022"}`)

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", payload)
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}
	})
}
