package api

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/controller"
	"github.com/ashalfarhan/realworld/api/dto"
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

	apiRoute.HandleFunc("/users", m.WithValidator(uc.RegisterUser, new(dto.RegisterUserDto))).Methods(http.MethodPost)
	apiRoute.HandleFunc("/users/login", m.WithValidator(uc.LoginUser, new(dto.LoginUserDto))).Methods(http.MethodPost)

	// User
	apiRoute.HandleFunc("/user", m.WithUser(uc.GetCurrentUser)).Methods(http.MethodGet)
	apiRoute.HandleFunc("/user", m.WithUser(uc.UpdateCurrentUser)).Methods(http.MethodPut)

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
	articleRoute.HandleFunc("/feed", m.WithUser(ac.GetFeed)).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}", ac.GetArticleBySlug).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}", m.WithUser(ac.DeleteArticle)).Methods(http.MethodDelete)
	articleRoute.HandleFunc("/{slug}", m.WithUser(ac.UpdateArticle)).Methods(http.MethodPut)
	articleRoute.HandleFunc("/{slug}/favorite", m.WithUser(ac.FavoriteArticle)).Methods(http.MethodPost)
	articleRoute.HandleFunc("/{slug}/favorite", m.WithUser(ac.UnFavoriteArticle)).Methods(http.MethodDelete)
	articleRoute.HandleFunc("/{slug}/comments", ac.GetArticleComments).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}/comments", m.WithUser(ac.CreateComment)).Methods(http.MethodPost)
	articleRoute.HandleFunc("/{slug}/comments/{id}", m.WithUser(ac.DeleteComment)).Methods(http.MethodDelete)

	return r
}
