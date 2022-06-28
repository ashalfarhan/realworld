package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/golang-jwt/jwt"
)

const (
	jwtSecret = "super-secret"
	jwtExp    = 20 * time.Minute
)

func GenerateJWT(u *model.User) (string, error) {
	now := time.Now()
	c := &jwt.StandardClaims{
		ExpiresAt: now.Add(jwtExp).Unix(),
		// Audience:  "client.com",
		Subject:  u.Username,
		IssuedAt: now.Unix(),
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
	return claim, nil
}
