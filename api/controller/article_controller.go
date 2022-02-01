package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
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
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	iu := c.authService.GetUserFromCtx(r)
	a, err := c.articleService.Create(r.Context(), d, iu.UserID)
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

	a, err := c.articleService.GetOneBySlug(r.Context(), slug)
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

	if err := c.articleService.DeleteArticle(r.Context(), slug, iu.UserID); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.articleService.GetAllTags(r.Context())
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"tags": tags,
	})
}

func (c *ArticleController) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	var d *dto.UpdateArticleDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	slug := mux.Vars(r)["slug"]
	iu := c.authService.GetUserFromCtx(r)

	ar, err := c.articleService.UpdateOneBySlug(r.Context(), iu.UserID, slug, d)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, response.M{
		"article": ar,
	})
}

func (c *ArticleController) GetFiltered(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var limit, offset int
	var err error
	l, o := q.Get("limit"), q.Get("offset")

	if l == "" {
		limit = 5
	} else {
		if limit, err = strconv.Atoi(l); err != nil {
			response.ClientError(w, err)
			return
		}
	}
	if o == "" {
		offset = 0
	} else {
		if offset, err = strconv.Atoi(o); err != nil {
			response.ClientError(w, err)
			return
		}
	}

	args := &conduit.ArticleArgs{
		Limit:  limit,
		Offset: offset,
	}

	v := validator.New()
	if err := v.Struct(args); err != nil {
		response.EntityError(w, err)
		return
	}

	articles, serr := c.articleService.GetArticles(r.Context(), args)
	if err != nil {
		response.Error(w, serr.Code, serr.Error)
		return
	}

	response.Ok(w, response.M{
		"articles": articles,
	})
}
