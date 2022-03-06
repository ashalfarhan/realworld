package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type userContextKey string

type AuthService struct {
	userCtxKey  userContextKey
	userService *UserService
	logger      *logrus.Entry
}

func NewAuthService(us *UserService) *AuthService {
	return &AuthService{"incoming-user", us, conduit.NewLogger("service", "AuthService")}
}

const (
	jwtSecret = "super-secret"
	jwtExp    = 20 * time.Hour
)

func (s AuthService) GenerateJWT(u *model.User) (string, error) {
	c := &conduit.ConduitClaims{
		UserID:   u.ID,
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExp).Unix(),
			Audience:  "client.com",
		},
	}

	s.logger.Infof("Generating jwt %#v", c)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	str, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		s.logger.Printf("Cannot sign jwt for: %#v, Reason: %v", u, err)
		return "", fmt.Errorf("cannot sign jwt: %w", err)
	}

	return str, nil
}

func (s AuthService) ParseJWT(str string) (*conduit.ConduitClaims, *ServiceError) {
	t, err := jwt.ParseWithClaims(str, new(conduit.ConduitClaims), getKey)
	if err != nil {
		s.logger.Printf("Cannot parse jwt for: %s, Reason: %v", str, err)
		return nil, CreateServiceError(http.StatusUnauthorized, fmt.Errorf("cannot parse jwt: %w", err))
	}

	claim, ok := t.Claims.(*conduit.ConduitClaims)
	if !ok {
		return nil, CreateServiceError(http.StatusUnauthorized, ErrInvalidClaim)
	}

	return claim, nil
}

func (s AuthService) GetUserFromCtx(r *http.Request) (*conduit.ConduitClaims, bool) {
	u, ok := r.Context().Value(s.userCtxKey).(*conduit.ConduitClaims)
	return u, ok
}

func (s AuthService) CreateUserCtx(parentCtx context.Context, claim *conduit.ConduitClaims) context.Context {
	return context.WithValue(parentCtx, s.userCtxKey, claim)
}

func (s AuthService) GetUserIDFromReq(r *http.Request) (string, *ServiceError) {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
	if len(authHeader) != 2 {
		return "", nil
	}

	jwt := authHeader[1]
	claim, err := s.ParseJWT(jwt)
	if err != nil {
		return "", err
	}

	return claim.UserID, nil
}

func (s AuthService) GetToken(r *http.Request) string {
	authHeader := strings.Split(r.Header.Get("Authorization"), "Token ")
	if len(authHeader) != 2 {
		return ""
	}

	return authHeader[1]
}

func getKey(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return []byte(jwtSecret), nil
}

func (s *AuthService) Login(ctx context.Context, d *dto.LoginUserFields) (*conduit.UserResponse, *ServiceError) {
	s.logger.Infof("POST Login %#v", d)
	u, sErr := s.userService.GetOne(ctx, &repository.FindOneUserFilter{
		Email:    d.Email,
		Username: d.Username,
	})
	if sErr != nil {
		return nil, sErr
	}

	if valid := u.ValidatePassword(d.Password); !valid {
		return nil, CreateServiceError(http.StatusBadRequest, ErrInvalidIdentity)
	}

	token, err := s.GenerateJWT(u)
	if err != nil {
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	res := &conduit.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
	}

	return res, nil
}

func (s *AuthService) Register(ctx context.Context, d *dto.RegisterUserFields) (*conduit.UserResponse, *ServiceError) {
	s.logger.Infof("POST Register %#v", d)
	u, sErr := s.userService.Insert(ctx, &RegisterArgs{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
	})

	if sErr != nil {
		return nil, sErr
	}

	token, err := s.GenerateJWT(u)
	if err != nil {
		return nil, CreateServiceError(http.StatusInternalServerError, nil)
	}

	res := &conduit.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
	}

	return res, nil
}
