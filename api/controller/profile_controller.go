package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils/jwt"
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
	iu := jwt.CurrentUser(r)
	profile, err := c.userService.FollowUser(r.Context(), iu, mux.Vars(r)["username"])
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"profile": profile,
	})
}

func (c *ProfileController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	iu := jwt.CurrentUser(r)
	profile, err := c.userService.UnfollowUser(r.Context(), iu, mux.Vars(r)["username"])
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"profile": profile,
	})
}

func (c *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	iu := jwt.CurrentUser(r)
	profile, err := c.userService.GetProfile(r.Context(), mux.Vars(r)["username"], iu)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.Ok(w, response.M{
		"profile": profile,
	})
}
