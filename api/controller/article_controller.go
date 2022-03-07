package controller

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/ashalfarhan/realworld/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ArticleController struct {
	articleService *service.ArticleService
	authService    *service.AuthService
}

func NewArticleController(s *service.Service) *ArticleController {
	return &ArticleController{s.ArticleService, s.AuthService}
}

func (c *ArticleController) CreateArticle(w http.ResponseWriter, r *http.Request) {
	d := r.Context().Value(dto.DtoCtxKey).(*dto.CreateArticleDto)

	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	a, err := c.articleService.CreateArticle(r.Context(), d.Article, iu.UserID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Created(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]

	uid, err := c.authService.GetUserIDFromReq(r)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	a, err := c.articleService.GetArticleBySlug(r.Context(), uid, slug)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user

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
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user

	ar, err := c.articleService.UpdateArticleBySlug(r.Context(), iu.UserID, slug, d)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, response.M{
		"article": ar.Serialize(),
	})
}

func (c *ArticleController) GetFiltered(w http.ResponseWriter, r *http.Request) {
	args, err := getArticleQueryParams(r.URL.Query())
	if err != nil {
		response.ClientError(w, err)
		return
	}

	id, serr := c.authService.GetUserIDFromReq(r)
	if err != nil {
		response.Error(w, serr.Code, serr.Error)
		return
	}
	args.UserID = id

	v := validator.New()
	if err := v.Struct(args); err != nil {
		response.EntityError(w, err)
		return
	}

	articles, serr := c.articleService.GetArticles(r.Context(), args)
	if serr != nil {
		response.Error(w, serr.Code, serr.Error)
		return
	}

	response.Ok(w, response.M{
		"articles":      articles.Serialize(),
		"articlesCount": len(articles),
	})
}

func (c *ArticleController) GetFeed(w http.ResponseWriter, r *http.Request) {
	args, err := getArticleQueryParams(r.URL.Query())
	if err != nil {
		response.ClientError(w, err)
		return
	}
	iu, _ := c.authService.GetUserFromCtx(r)
	args.UserID = iu.UserID

	v := validator.New()
	if err := v.Struct(args); err != nil {
		response.EntityError(w, err)
		return
	}

	articles, serr := c.articleService.GetArticlesFeed(r.Context(), args)
	if serr != nil {
		response.Error(w, serr.Code, serr.Error)
		return
	}

	response.Ok(w, response.M{
		"articles":      articles.Serialize(),
		"articlesCount": len(articles),
	})
}

func (c *ArticleController) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	a, err := c.articleService.FavoriteArticleBySlug(r.Context(), iu.UserID, slug)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) UnFavoriteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu, _ := c.authService.GetUserFromCtx(r) // There will always be a user
	a, err := c.articleService.UnfavoriteArticleBySlug(r.Context(), iu.UserID, slug)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, response.M{
		"article": a.Serialize(),
	})
}

func getArticleQueryParams(q url.Values) (*repository.FindArticlesArgs, error) {
	var limit, offset int
	var err error
	l, o := q.Get("limit"), q.Get("offset")

	if l == "" {
		limit = 5
	} else {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return nil, err
		}
	}

	if o == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(o)
		if err != nil {
			return nil, err
		}
	}

	args := &repository.FindArticlesArgs{
		Limit:  limit,
		Offset: offset,
	}

	return args, nil
}
