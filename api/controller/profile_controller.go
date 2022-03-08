package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/mux"
)

type ProfileController struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewProfileController(s *service.Service) *ProfileController {
	return &ProfileController{s.UserService, s.AuthService}
}

func (c *ProfileController) FollowUser(w http.ResponseWriter, r *http.Request) {
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	profile, err := c.userService.FollowUser(r.Context(), iu.UserID, mux.Vars(r)["username"])
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"profile": profile,
	})
}

func (c *ProfileController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	profile, err := c.userService.UnfollowUser(r.Context(), iu.UserID, mux.Vars(r)["username"])
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"profile": profile,
	})
}

func (c *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	profile, err := c.userService.GetProfile(r.Context(), mux.Vars(r)["username"], iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}
	response.Ok(w, response.M{
		"profile": profile,
	})
}
