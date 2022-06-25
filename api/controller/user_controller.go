package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/ashalfarhan/realworld/utils/jwt"
)

type UserController struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewUserController(s *service.Service) *UserController {
	return &UserController{s.UserService, s.AuthService}
}

func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	iu, _ := jwt.CurrentUser(r) // There will always be a user

	u, err := c.userService.GetOneById(r.Context(), iu.Subject)
	if err != nil {
		response.Err(w, err)
		return
	}

	res := &model.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    jwt.GetToken(r),
	}
	response.Ok(w, response.M{
		"user": res,
	})
}

func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	req := new(model.UpdateUserDto)
	err := utils.ValidateDTO(r, req)
	if err != nil {
		response.Err(w, err)
		return
	}

	iu, _ := jwt.CurrentUser(r) // There will always be a user
	u, err := c.userService.Update(r.Context(), req.User, iu.Subject)
	if err != nil {
		response.Err(w, err)
		return
	}

	res := &model.UserResponse{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    jwt.GetToken(r),
	}
	response.Accepted(w, response.M{
		"user": res,
	})
}
