package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/service"
)

type UserController struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewUserController(s *service.Service) *UserController {
	return &UserController{s.UserService, s.AuthService}
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	d := r.Context().Value(dto.DtoCtxKey).(*dto.RegisterUserDto)
	res, err := c.authService.Register(r.Context(), d.User)

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Created(w, response.M{
		"user": res,
	})
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	d := r.Context().Value(dto.DtoCtxKey).(*dto.LoginUserDto)
	res, err := c.authService.Login(r.Context(), d.User)

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"user": res,
	})
}

func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user

	u, err := c.userService.GetOneById(r.Context(), iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	res := &conduit.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
	token := c.authService.GetToken(r)

	if token != "" {
		res.Token = token
	}

	response.Ok(w, response.M{
		"user": res,
	})
}

func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	d := r.Context().Value(dto.DtoCtxKey).(*dto.UpdateUserDto)

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	u, err := c.userService.Update(r.Context(), d, iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	res := &conduit.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
	}
	token := c.authService.GetToken(r)
	if token != "" {
		res.Token = token
	}

	response.Accepted(w, response.M{
		"user": res,
	})
}
