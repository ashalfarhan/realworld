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
		userService: s.US,
		authService: s.AS,
	}
}

func (p *ProfileController) FollowUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]

	iu := p.authService.GetUserFromCtx(r)

	if err := p.userService.FollowUser(iu.UserID, uname); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Success(w, http.StatusAccepted, nil)
}

func (p *ProfileController) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]

	iu := p.authService.GetUserFromCtx(r)

	if err := p.userService.UnfollowUser(iu.UserID, uname); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Success(w, http.StatusAccepted, nil)
}

func (p *ProfileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	uname := mux.Vars(r)["username"]
	u, err := p.userService.GetOne(&dto.LoginUserDto{Username: uname})
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	iu := p.authService.GetUserFromCtx(r)
	following := p.userService.IsFollowing(iu.UserID, u.ID)
	res := &conduit.ProfileResponse{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}

	response.Success(w, http.StatusOK, res)
}
