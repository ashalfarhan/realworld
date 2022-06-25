package middleware

import (
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/utils/jwt"
)

func (m *ConduitMiddleware) WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
		if len(authHeader) != 2 {
			response.UnauthorizeError(w, "No token")
			return
		}

		token := authHeader[1]
		claim, err := jwt.ParseJWT(token)
		if err != nil {
			response.Err(w, err)
			return
		}

		ctx := jwt.CreateUserCtx(r.Context(), claim)
		next(w, r.WithContext(ctx))
	}
}
