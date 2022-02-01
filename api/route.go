package api

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/controller"
	"github.com/ashalfarhan/realworld/api/middleware"
	"github.com/ashalfarhan/realworld/service"
	"github.com/gorilla/mux"
)

func InitRoutes(s *service.Service) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controller.Hello).Methods(http.MethodGet)
	apiRoute := r.PathPrefix("/api").Subrouter()
	m := middleware.NewMiddleware(s)

	// Users
	uc := controller.NewUserController(s)
	usersRoute := apiRoute.PathPrefix("/users").Subrouter()

	usersRoute.HandleFunc("/register", uc.RegisterUser).Methods(http.MethodPost)
	usersRoute.HandleFunc("/login", uc.LoginUser).Methods(http.MethodPost)

	// User
	userRoute := apiRoute.PathPrefix("/user").Subrouter()

	userRoute.HandleFunc("", m.WithUser(uc.GetCurrentUser)).Methods(http.MethodGet)
	userRoute.HandleFunc("", m.WithUser(uc.UpdateCurrentUser)).Methods(http.MethodPut)

	// Profile
	pc := controller.NewProfileController(s)
	profileRoute := apiRoute.PathPrefix("/profiles").Subrouter()

	profileRoute.HandleFunc("/{username}", m.WithUser(pc.GetProfile)).Methods(http.MethodGet)
	profileRoute.HandleFunc("/{username}/follow", m.WithUser(pc.FollowUser)).Methods(http.MethodPost)
	profileRoute.HandleFunc("/{username}/follow", m.WithUser(pc.UnfollowUser)).Methods(http.MethodDelete)

	// Article
	ac := controller.NewArticleController(s)
	articleRoute := apiRoute.PathPrefix("/articles").Subrouter()
	// Tags
	apiRoute.HandleFunc("/tags", ac.GetAllTags).Methods(http.MethodGet)

	articleRoute.HandleFunc("", ac.GetFiltered).Methods(http.MethodGet)
	articleRoute.HandleFunc("", m.WithUser(ac.CreateArticle)).Methods(http.MethodPost)
	articleRoute.HandleFunc("/{slug}", ac.GetArticleBySlug).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}", m.WithUser(ac.DeleteArticle)).Methods(http.MethodDelete)
	articleRoute.HandleFunc("/{slug}", m.WithUser(ac.UpdateArticle)).Methods(http.MethodPut)

	return r
}
