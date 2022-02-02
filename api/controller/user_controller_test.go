package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashalfarhan/realworld/api/response"
)

type DtoError map[string]map[string][]string

func TestLoginController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if not provide email or username", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"user": response.M{
				"password": "secret",
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.LoginUser(w, req)
		res := w.Result()

		defer res.Body.Close()
		var body DtoError
		json.NewDecoder(res.Body).Decode(&body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["username"]) == 0 && len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in username and email, but got %v", body["errors"])
		}
	})
}

func TestRegisterController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if not provide email and username", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"user": response.M{
				"password": "secret",
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		defer res.Body.Close()
		var body DtoError

		json.NewDecoder(res.Body).Decode(&body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["username"]) == 0 && len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in username and email, but got %v", body["errors"])
		}
	})

	t.Run("Should response error if not password is less than 8", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"user": response.M{
				"email":    "johndoe@mail.com",
				"username": "doejohn",
				"password": "secret",
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		defer res.Body.Close()
		var body DtoError

		json.NewDecoder(res.Body).Decode(&body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if body["errors"]["password"][0] == "min 8 " {
			t.Fatalf("expected validation errors in password, but got %v", body["errors"])
		}
	})

	t.Run("Should response error if not a valid email", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"user": response.M{
				"email":    "johndoe@mail@~asdasd.com",
				"username": "doejohn",
				"password": "secret2022",
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.RegisterUser(w, req)
		res := w.Result()

		defer res.Body.Close()
		var body DtoError

		json.NewDecoder(res.Body).Decode(&body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if len(body["errors"]["email"]) == 0 {
			t.Fatalf("expected validation errors in email, but got %v", body["errors"])
		}
	})
}

func TestUpdateUserController(t *testing.T) {
	userController := &UserController{}
	t.Run("Should response error if new password length is not the min (8)", func(t *testing.T) {
		b, _ := json.Marshal(response.M{
			"user": response.M{
				"password": "asd",
			},
		})

		req := httptest.NewRequest(http.MethodPut, "/api/user", bytes.NewReader(b))
		w := httptest.NewRecorder()

		userController.UpdateCurrentUser(w, req)
		res := w.Result()

		defer res.Body.Close()
		var body DtoError

		json.NewDecoder(res.Body).Decode(&body)

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected StatusCode to be %d, but got %v", http.StatusUnprocessableEntity, res.StatusCode)
		}

		if body["errors"]["password"][0] == "min 8 " {
			t.Fatalf("expected validation errors in password, but got %v", body["errors"])
		}
	})
}
