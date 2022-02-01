package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/golang-jwt/jwt"
	
)

type userContextKey string

type AuthService struct {
	userCtxKey userContextKey
}

func NewAuthService() *AuthService {
	return &AuthService{
		userCtxKey: "incoming-user",
	}
}

const (
	jwtSecret = "super-secret"
	jwtExp    = 20 * time.Hour
)

func (AuthService) GenerateJWT(u *model.User) (string, error) {
	c := &conduit.ConduitClaims{
		UserID:   u.ID,
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExp).Unix(),
			Audience:  "client.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	str, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("cannot sign jwt: %w", err)
	}

	return str, nil
}

func (AuthService) ParseJWT(str string) (*conduit.ConduitClaims, error) {
	t, err := jwt.ParseWithClaims(str, &conduit.ConduitClaims{}, getKey)
	if err != nil {
		return nil, fmt.Errorf("cannot parse jwt: %w", err)
	}

	claim, ok := t.Claims.(*conduit.ConduitClaims)
	if !ok {
		return nil, errors.New("invalid claim")
	}

	return claim, nil
}

func (a AuthService) GetUserFromCtx(r *http.Request) *conduit.ConduitClaims {
	return r.Context().Value(a.userCtxKey).(*conduit.ConduitClaims)
}

func (a AuthService) CreateUserCtx(parentCtx context.Context, claim *conduit.ConduitClaims) context.Context {
	return context.WithValue(parentCtx, a.userCtxKey, claim)
}

func getKey(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return []byte(jwtSecret), nil
}
