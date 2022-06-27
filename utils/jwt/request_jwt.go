package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/ashalfarhan/realworld/model"
	"github.com/golang-jwt/jwt"
)

func CreateUserCtx(ctx context.Context, claim *jwt.StandardClaims) context.Context {
	return context.WithValue(ctx, userCtx, claim)
}

// Get User ID from request ctx (required auth endpoint).
// Return empty string if no user from the ctx
func CurrentUser(r *http.Request) string {
	u, ok := r.Context().Value(userCtx).(*jwt.StandardClaims)
	if !ok {
		return ""
	}
	return u.Subject
}

// Get User ID from request.
// Used for non-auth endpoint to retrieve user id (empty string if no token).
// Error returned will be if invalid jwt
func GetUsernameFromReq(r *http.Request) (string, *model.ConduitError) {
	token := GetToken(r)
	if token == "" {
		return token, nil
	}
	claim, err := ParseJWT(token)
	if err != nil {
		return "", err
	}
	return claim.Subject, nil
}

// Get JWT Token from Request Header
func GetToken(r *http.Request) string {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
	if len(authHeader) != 2 {
		return ""
	}
	return authHeader[1]
}
