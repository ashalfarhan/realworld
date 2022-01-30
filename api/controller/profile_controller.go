package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
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

	iu := c.authService.GetUserFromCtx(r)

	if err := c.userService.FollowUser(iu.UserID, uname); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Success(w, http.StatusAccepted, nil)
}

func (c *ProfileController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]

	iu := c.authService.GetUserFromCtx(r)

	if err := c.userService.UnfollowUser(iu.UserID, uname); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Success(w, http.StatusAccepted, nil)
}

func (c *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]
	u, err := c.userService.GetOne(&dto.LoginUserDto{Username: uname})
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	iu := c.authService.GetUserFromCtx(r)
	following := c.userService.IsFollowing(iu.UserID, u.ID)
	res := &conduit.ProfileResponse{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}

	response.Success(w, http.StatusOK, res)
}
