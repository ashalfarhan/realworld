package service

import (
	"context"
	"net/http"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/utils/jwt"
	"github.com/ashalfarhan/realworld/utils/logger"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(us *UserService) *AuthService {
	return &AuthService{
		userService: us,
	}
}

func (s AuthService) Login(ctx context.Context, d *model.LoginUserFields) (*model.UserResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST Login %#v", d)
	u, sErr := s.userService.GetOne(ctx, &repository.FindOneUserFilter{
		Email:    d.Email,
		Username: d.Username,
	})
	if sErr != nil {
		return nil, sErr
	}

	if valid := u.ValidatePassword(d.Password); !valid {
		return nil, conduit.BuildError(http.StatusBadRequest, ErrInvalidIdentity)
	}

	token, err := jwt.GenerateJWT(u)
	if err != nil {
		return nil, conduit.GeneralError
	}

	res := &model.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
	}

	return res, nil
}

func (s AuthService) Register(ctx context.Context, d *model.RegisterUserFields) (*model.UserResponse, *model.ConduitError) {
	log := logger.GetCtx(ctx)
	log.Infof("POST Register %#v", *d)
	u, sErr := s.userService.Insert(ctx, d)
	if sErr != nil {
		return nil, sErr
	}

	token, err := jwt.GenerateJWT(u)
	if err != nil {
		return nil, conduit.GeneralError
	}

	res := &model.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
	}

	return res, nil
}
