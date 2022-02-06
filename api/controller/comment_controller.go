package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/db/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func (c *ArticleController) CreateComment(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	var d *dto.CreateCommentDto
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		response.ClientError(w, err)
		return
	}

	v := validator.New()
	if err := v.Struct(d); err != nil {
		response.EntityError(w, err)
		return
	}

	iu, _ := c.authService.GetUserFromCtx(r)
	d.AuthorID = iu.UserID
	d.ArticleSlug = slug
	comm, err := c.articleService.CreateComment(r.Context(), d)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Created(w, response.M{
		"comment": comm,
	})
}

func (c *ArticleController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, id := params["slug"], params["id"]
	iu, _ := c.authService.GetUserFromCtx(r)
	if err := c.articleService.DeleteCommentByID(r.Context(), id, iu.UserID); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetArticleComments(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var limit, offset int
	var err error
	l, o := q.Get("limit"), q.Get("offset")

	if l == "" {
		limit = 5
	} else {
		limit, err = strconv.Atoi(l)
		if err != nil {
			response.ClientError(w, err)
			return
		}
	}

	if o == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(o)
		if err != nil {
			response.ClientError(w, err)
			return
		}
	}

	userID, sErr := c.authService.GetUserIDFromReq(r)
	if sErr != nil {
		response.Error(w, sErr.Code, sErr.Error)
		return
	}

	slug := mux.Vars(r)["slug"]
	a, sErr := c.articleService.GetArticleBySlug(r.Context(), userID, slug)
	if sErr != nil {
		response.Error(w, sErr.Code, sErr.Error)
		return
	}

	args := &repository.FindCommentsByArticleIDArgs{
		Limit:     limit,
		Offset:    offset,
		ArticleID: a.ID,
	}

	v := validator.New()
	if err := v.Struct(args); err != nil {
		response.EntityError(w, err)
		return
	}

	comms, sErr := c.articleService.GetComments(r.Context(), args, slug)
	if sErr != nil {
		response.Error(w, sErr.Code, sErr.Error)
		return
	}

	response.Ok(w, response.M{
		"comments": comms,
	})
}
