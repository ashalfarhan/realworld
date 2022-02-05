package middleware

import (
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/api/response"
)

func (m *ConduitMiddleware) WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
		if len(authHeader) != 2 {
			response.UnauthorizeError(w, "No token")
			return
		}

		jwt := authHeader[1]
		claim, err := m.authService.ParseJWT(jwt)
		if err != nil {
			response.UnauthorizeError(w, err.Error())
			return
		}

		ctx := m.authService.CreateUserCtx(r.Context(), claim)
		next(w, r.WithContext(ctx))
	}
}
