package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/api/response"
	"github.com/gorilla/mux"
)

func (c *ArticleController) CreateComment(w http.ResponseWriter, r *http.Request) {
	d := r.Context().Value(dto.DtoCtxKey).(*dto.CreateCommentDto)
	iu, _ := c.authService.GetUserFromCtx(r)

	d.AuthorID = iu.UserID
	d.ArticleSlug = mux.Vars(r)["slug"]

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
	iu, _ := c.authService.GetUserFromCtx(r)
	if err := c.articleService.DeleteCommentByID(r.Context(), mux.Vars(r)["id"], iu.UserID); err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetArticleComments(w http.ResponseWriter, r *http.Request) {
	userID, err := c.authService.GetUserIDFromReq(r)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	comms, err := c.articleService.GetComments(r.Context(), mux.Vars(r)["slug"], userID)
	if err != nil {
		response.Error(w, err.Code, err.Error)
		return
	}

	response.Ok(w, response.M{
		"comments": comms,
	})
}
