package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/golang-jwt/jwt"
)

type UserCtxKey string

var userCtx UserCtxKey = "incoming-user"

const (
	jwtSecret = "super-secret"
	jwtExp    = 20 * time.Minute
)

func GenerateJWT(u *model.User) (string, error) {
	now := time.Now()
	c := &jwt.StandardClaims{
		ExpiresAt: now.Add(jwtExp).Unix(),
		Audience:  "client.com", // TODO: Change this
		Subject:   u.ID,
		IssuedAt:  now.Unix(),
		// Issuer:    u.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	str, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("cannot sign jwt: %w", err)
	}
	return str, nil
}

func ParseJWT(str string) (*jwt.StandardClaims, *model.ConduitError) {
	t, err := jwt.ParseWithClaims(str, new(jwt.StandardClaims), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, conduit.BuildError(401, fmt.Errorf("cannot parse jwt: %w", err))
	}
	claim, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, conduit.BuildError(401, errors.New("invalid claim"))
	}
	if err = claim.Valid(); err != nil {
		return nil, conduit.BuildError(401, err)
	}
	return claim, nil
}

func CurrentUser(r *http.Request) (*jwt.StandardClaims, bool) {
	u, ok := r.Context().Value(userCtx).(*jwt.StandardClaims)
	return u, ok
}

func CreateUserCtx(ctx context.Context, claim *jwt.StandardClaims) context.Context {
	return context.WithValue(ctx, userCtx, claim)
}

func GetUserIDFromReq(r *http.Request) (string, *model.ConduitError) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
	if len(authHeader) != 2 {
		return "", nil
	}

	jwt := authHeader[1]
	claim, err := ParseJWT(jwt)
	if err != nil {
		return "", err
	}

	return claim.Subject, nil
}

func GetToken(r *http.Request) string {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
	if len(authHeader) != 2 {
		return ""
	}
	return authHeader[1]
}
