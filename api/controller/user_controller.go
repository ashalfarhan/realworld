package controller

import (
	"encoding/json"
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

func NewUserController(svc *service.Service) *UserController {
	return &UserController{
		userService: svc.US,
		authService: svc.AS,
	}
}

func (c *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.RegisterUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	if err := dto.ValidateDto(d, v); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	u, err := c.userService.CreateOne(d)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, map[string]interface{}{
		"user": u,
	})
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.LoginUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	if err := dto.ValidateDto(d, v); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	u, err := c.userService.GetOne(d)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid identity or password")
		return
	}

	if valid := u.ValidatePassword(d.Password); !valid {
		response.Error(w, http.StatusBadRequest, "Invalid identity or password")
		return
	}

	token, err := c.authService.GenerateJWT(u)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := &conduit.UserAuthResponse{
		Email:    u.Email,
		Username: u.Username,
		Token:    token,
		Bio:      u.Bio,
		Image:    u.Image,
	}
	response.Success(w, http.StatusOK, res)
}

func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	iu := c.authService.GetUserFromCtx(r)
	u, err := c.userService.GetOneById(iu.UserID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"user": u,
	})
}

func (c *UserController) UpdateCurrentUser(w http.ResponseWriter, r *http.Request) {
	var d *dto.UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New()
	if err := dto.ValidateDto(d, v); err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	iu := c.authService.GetUserFromCtx(r)
	if err := c.userService.Update(d, iu.UserID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusOK, nil)
}
