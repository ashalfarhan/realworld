package api

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/controller"
	"github.com/ashalfarhan/realworld/api/middleware"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/mux"
)

func InitRoutes(s *service.Service) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", controller.Hello).Methods(http.MethodGet)
	apiRoute := r.PathPrefix("/api").Subrouter()

	// User
	uc := controller.NewUserController(s)
	usersRoute := apiRoute.PathPrefix("/users").Subrouter()
	usersRoute.HandleFunc("/register", uc.RegisterUser).Methods(http.MethodPost)
	usersRoute.HandleFunc("/login", uc.LoginUser).Methods(http.MethodPost)

	userRoute := apiRoute.PathPrefix("/user").Subrouter()
	userRoute.HandleFunc("", middleware.WithUser(s.AS, uc.GetCurrentUser)).Methods(http.MethodGet)
	userRoute.HandleFunc("", middleware.WithUser(s.AS, uc.UpdateCurrentUser)).Methods(http.MethodPut)

	// Profile
	pc := controller.NewProfileController(s)
	profileRoute := apiRoute.PathPrefix("/profiles").Subrouter()
	profileRoute.HandleFunc("/{username}", middleware.WithUser(s.AS, pc.GetProfile)).Methods(http.MethodGet)
	profileRoute.HandleFunc("/{username}/follow", middleware.WithUser(s.AS, pc.FollowUser)).Methods(http.MethodPost)
	profileRoute.HandleFunc("/{username}/follow", middleware.WithUser(s.AS, pc.UnfollowUser)).Methods(http.MethodDelete)

	return r
}
