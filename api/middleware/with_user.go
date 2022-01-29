package middleware

import (
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/service"
)

func WithUser(authService *service.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			response.Error(w, http.StatusUnauthorized, "No jwt provided")
			return
		}

		jwt := authHeader[1]
		claim, err := authService.ParseJWT(jwt)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := authService.CreateUserCtx(r.Context(), claim)
		next(w, r.WithContext(ctx))
	})
}
