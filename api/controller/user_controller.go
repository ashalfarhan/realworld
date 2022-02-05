package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/service"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewUserController(s *service.Service) *UserController {
	return &UserController{
		userService: s.UserService,
		authService: s.AuthService,
	}
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.RegisterUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	u, err := c.userService.Register(r.Context(), &service.RegisterArgs{
		Email:    d.User.Email,
		Username: d.User.Username,
		Password: d.User.Password,
	})

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Created(w, response.M{
		"user": u,
	})
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.LoginUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	u, err := c.userService.GetOne(r.Context(), &service.GetOneArgs{
		Email:    d.User.Email,
		Username: d.User.Username,
	})

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	if valid := u.ValidatePassword(d.User.Password); !valid {
		response.ClientError(w, errors.New("invalid identity or password"))
		return
	}

	token, serr := c.authService.GenerateJWT(u)
	if serr != nil {
		response.InternalError(w)
		return
	}

	res := &conduit.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
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

	response.Ok(w, response.M{
		"user": u,
	})
}

func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	u, err := c.userService.Update(r.Context(), d, iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, response.M{
		"user": u,
	})
}
