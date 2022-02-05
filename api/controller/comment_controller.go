package controller

import (
	"encoding/json"
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
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
	d.Comment.AuthorID = iu.UserID
	d.Comment.ArticleSlug = slug
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
