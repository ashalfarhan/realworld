package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/conduit"
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
	if err := utils.ValidateDTO(r, req); err != nil {
		response.Err(w, err)
		return
	}

	iu := jwt.CurrentUser(r)
	a, err := c.articleService.CreateArticle(r.Context(), req.Article, iu)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Created(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	uid, err := jwt.GetUsernameFromReq(r)
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
	iu := jwt.CurrentUser(r)
	if err := c.articleService.DeleteArticle(r.Context(), mux.Vars(r)["slug"], iu); err != nil {
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
	if err := utils.ValidateDTO(r, req); err != nil {
		response.Err(w, err)
		return
	}

	iu := jwt.CurrentUser(r)
	ar, err := c.articleService.UpdateArticleBySlug(r.Context(), iu, mux.Vars(r)["slug"], req.Article)
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
		response.Err(w, err)
		return
	}

	articles, err := c.articleService.GetArticles(r.Context(), args)
	if err != nil {
		response.Err(w, err)
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
		response.Err(w, err)
		return
	}

	args.Username = jwt.CurrentUser(r)
	articles, err := c.articleService.GetArticlesFeed(r.Context(), args)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"articles":      articles.Serialize(),
		"articlesCount": len(articles),
	})
}

func (c *ArticleController) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	iu := jwt.CurrentUser(r)
	a, err := c.articleService.FavoriteArticleBySlug(r.Context(), iu, mux.Vars(r)["slug"])
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, response.M{
		"article": a.Serialize(),
	})
}

func (c *ArticleController) UnFavoriteArticle(w http.ResponseWriter, r *http.Request) {
	iu := jwt.CurrentUser(r) // There will always be a user
	a, err := c.articleService.UnfavoriteArticleBySlug(r.Context(), iu, mux.Vars(r)["slug"])
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, response.M{
		"article": a.Serialize(),
	})
}

func getArticleQueryParams(q url.Values) (*repository.FindArticlesArgs, *model.ConduitError) {
	var err error
	limit, offset := q.Get("limit"), q.Get("offset")
	args := &repository.FindArticlesArgs{
		Tag:       q.Get("tag"),
		Author:    q.Get("author"),
		Favorited: q.Get("favorited"),
	}

	if limit == "" {
		// Default if not specified
		limit = "5"
	}
	if args.Limit, err = strconv.Atoi(limit); err != nil {
		return nil, conduit.BuildError(400, err)
	}

	if offset == "" {
		// Default if not specified
		offset = "0"
	}
	if args.Offset, err = strconv.Atoi(offset); err != nil {
		return nil, conduit.BuildError(400, err)
	}

	v := validator.New()
	if err = v.Struct(args); err != nil {
		return nil, conduit.BuildError(http.StatusUnprocessableEntity, err)
	}

	return args, nil
}
