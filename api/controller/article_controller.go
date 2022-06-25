package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/persistence/repository"
	"github.com/ashalfarhan/realworld/service"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/ashalfarhan/realworld/utils/jwt"
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
	req := new(model.CreateArticleDto)
	err := utils.ValidateDTO(r, req)
	if err != nil {
		response.Err(w, err)
		return
	}

	iu, _ := jwt.CurrentUser(r) // There will always be a user
	a, err := c.articleService.CreateArticle(r.Context(), req.Article, iu.Subject)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Created(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	uid, err := jwt.GetUserIDFromReq(r)
	if err != nil {
		response.Err(w, err)
		return
	}

	a, err := c.articleService.GetArticleBySlug(r.Context(), uid, mux.Vars(r)["slug"])
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	iu, _ := jwt.CurrentUser(r) // There will always be a user

	if err := c.articleService.DeleteArticle(r.Context(), mux.Vars(r)["slug"], iu.Subject); err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := c.articleService.GetAllTags(r.Context())
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"tags": tags,
	})
}

func (c *ArticleController) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	req := new(model.UpdateArticleDto)
	err := utils.ValidateDTO(r, req)
	if err != nil {
		response.Err(w, err)
		return
	}

	iu, _ := jwt.CurrentUser(r) // There will always be a user
	ar, err := c.articleService.UpdateArticleBySlug(r.Context(), iu.Subject, mux.Vars(r)["slug"], req.Article)
	if err != nil {
		response.Err(w, err)
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

	id, serr := jwt.GetUserIDFromReq(r)
	if err != nil {
		response.Err(w, serr)
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
		response.Err(w, serr)
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
	iu, _ := jwt.CurrentUser(r)
	args.UserID = iu.Subject

	v := validator.New()
	if err := v.Struct(args); err != nil {
		response.EntityError(w, err)
		return
	}

	articles, serr := c.articleService.GetArticlesFeed(r.Context(), args)
	if serr != nil {
		response.Err(w, serr)
		return
	}

	response.Ok(w, response.M{
		"articles":      articles.Serialize(),
		"articlesCount": len(articles),
	})
}

func (c *ArticleController) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu, _ := jwt.CurrentUser(r) // There will always be a user
	a, err := c.articleService.FavoriteArticleBySlug(r.Context(), iu.Subject, slug)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) UnFavoriteArticle(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	iu, _ := jwt.CurrentUser(r) // There will always be a user
	a, err := c.articleService.UnfavoriteArticleBySlug(r.Context(), iu.Subject, slug)
	if err != nil {
		response.Err(w, err)
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
