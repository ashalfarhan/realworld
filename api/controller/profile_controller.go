package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/mux"
)

type ProfileController struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewProfileController(s *service.Service) *ProfileController {
	return &ProfileController{
		userService: s.UserService,
		authService: s.AuthService,
	}
}

func (c *ProfileController) FollowUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	profile, err := c.userService.FollowUser(r.Context(), iu.UserID, uname)

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	res := &conduit.ProfileResponse{
		Username:  profile.Username,
		Bio:       profile.Bio,
		Image:     profile.Image,
		Following: true,
	}

	response.Ok(w, response.M{
		"profile": res,
	})
}

func (c *ProfileController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	profile, err := c.userService.UnfollowUser(r.Context(), iu.UserID, uname)

	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	res := &conduit.ProfileResponse{
		Username:  profile.Username,
		Bio:       profile.Bio,
		Image:     profile.Image,
		Following: false,
	}

	response.Ok(w, response.M{
		"profile": res,
	})
}

func (c *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]
	u, err := c.userService.GetOne(r.Context(), &repository.FindOneUserFilter{Username: uname})
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	following := c.userService.IsFollowing(r.Context(), iu.UserID, u.ID)

	res := &conduit.ProfileResponse{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}

	response.Ok(w, response.M{
		"profile": res,
	})
}
