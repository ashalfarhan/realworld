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
	r.Use(middleware.InjectReqID)

	r.HandleFunc("/", controller.Hello).Methods(http.MethodGet)
	apiRoute := r.PathPrefix("/api").Subrouter()

	// Auth
	auth := controller.NewAuthController(s)
	apiRoute.HandleFunc("/users", auth.RegisterUser).Methods(http.MethodPost)
	apiRoute.HandleFunc("/users/login", auth.LoginUser).Methods(http.MethodPost)

	// User
	uc := controller.NewUserController(s)
	apiRoute.HandleFunc("/user", middleware.WithUser(uc.GetCurrentUser)).Methods(http.MethodGet)
	apiRoute.HandleFunc("/user", middleware.WithUser(uc.UpdateCurrentUser)).Methods(http.MethodPut)

	// Profile
	pc := controller.NewProfileController(s)
	profileRoute := apiRoute.PathPrefix("/profiles").Subrouter()
	profileRoute.HandleFunc("/{username}", middleware.WithUser(pc.GetProfile)).Methods(http.MethodGet)
	profileRoute.HandleFunc("/{username}/follow", middleware.WithUser(pc.FollowUser)).Methods(http.MethodPost)
	profileRoute.HandleFunc("/{username}/follow", middleware.WithUser(pc.UnfollowUser)).Methods(http.MethodDelete)

	// Article
	ac := controller.NewArticleController(s)
	apiRoute.HandleFunc("/tags", ac.GetAllTags).Methods(http.MethodGet)
	apiRoute.HandleFunc("/articles", ac.GetFiltered).Methods(http.MethodGet)
	apiRoute.HandleFunc("/articles", middleware.WithUser(ac.CreateArticle)).Methods(http.MethodPost)
	articleRoute := apiRoute.PathPrefix("/articles").Subrouter()
	articleRoute.HandleFunc("/feed", middleware.WithUser(ac.GetFeed)).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}", ac.GetArticleBySlug).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}", middleware.WithUser(ac.DeleteArticle)).Methods(http.MethodDelete)
	articleRoute.HandleFunc("/{slug}", middleware.WithUser(ac.UpdateArticle)).Methods(http.MethodPut)
	articleRoute.HandleFunc("/{slug}/favorite", middleware.WithUser(ac.FavoriteArticle)).Methods(http.MethodPost)
	articleRoute.HandleFunc("/{slug}/favorite", middleware.WithUser(ac.UnFavoriteArticle)).Methods(http.MethodDelete)
	articleRoute.HandleFunc("/{slug}/comments", ac.GetArticleComments).Methods(http.MethodGet)
	articleRoute.HandleFunc("/{slug}/comments", middleware.WithUser(ac.CreateComment)).Methods(http.MethodPost)
	articleRoute.HandleFunc("/{slug}/comments/{id}", middleware.WithUser(ac.DeleteComment)).Methods(http.MethodDelete)

	return r
}
