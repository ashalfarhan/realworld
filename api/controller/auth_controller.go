package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(s *service.Service) *AuthController {
	return &AuthController{
		service: s.AuthService,
	}
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	req := new(model.RegisterUserDto)
	if err := utils.ValidateDTO(r, req); err != nil {
		response.Err(w, err)
		return
	}

	res, err := c.service.Register(r.Context(), req.User)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Created(w, response.M{
		"user": res,
	})
}

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	req := new(model.LoginUserDto)
	if err := utils.ValidateDTO(r, req); err != nil {
		response.Err(w, err)
		return
	}

	res, err := c.service.Login(r.Context(), req.User)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"user": res,
	})
}
