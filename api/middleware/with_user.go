package middleware

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/utils/jwt"
)

func WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := jwt.GetToken(r)
		if token == "" {
			response.UnauthorizeError(w, "No token")
			return
		}
		claim, err := jwt.ParseJWT(token)
		if err != nil {
			response.Err(w, err)
			return
		}
		ctx := jwt.CreateUserCtx(r.Context(), claim)
		next(w, r.WithContext(ctx))
	}
}
