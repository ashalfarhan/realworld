package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ArticleController struct {
	articleService *service.ArticleService
	authService    *service.AuthService
}

func NewArticleController(s *service.Service) *ArticleController {
	return &ArticleController{
		articleService: s.ArticleService,
		authService:    s.AuthService,
	}
}

func (c *ArticleController) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var d *dto.CreateArticleDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := dto.ValidateDto(d, v); err != nil {
		response.EntityError(w, err)
		return
	}

	iu := c.authService.GetUserFromCtx(r)
	a, err := c.articleService.Create(d, iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Created(w, response.M{
		"article": a,
	})
}

func (c *ArticleController) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	a, err := c.articleService.GetOneBySlug(slug)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"article": a,
	})
}

func (c *ArticleController) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu := c.authService.GetUserFromCtx(r)

	if err := c.articleService.DeleteArticle(slug, iu.UserID); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w)
}

func (c *ArticleController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.articleService.GetAllTags()
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"tags": tags,
	})
}
