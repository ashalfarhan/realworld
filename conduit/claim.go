package conduit

import "github.com/golang-jwt/jwt"

type ConduitClaims struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}
