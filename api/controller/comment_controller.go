package controller

import (
	"net/http"

	"github.com/ashalfarhan/realworld/api/response"
	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils"
	"github.com/ashalfarhan/realworld/utils/jwt"
	"github.com/gorilla/mux"
)

func (c *ArticleController) CreateComment(w http.ResponseWriter, r *http.Request) {
	req := new(model.CreateCommentDto)
	if err := utils.ValidateDTO(r, req); err != nil {
		response.Err(w, err)
		return
	}
	req.AuthorUsername = jwt.CurrentUser(r)
	req.ArticleSlug = mux.Vars(r)["slug"]

	comm, err := c.articleService.CreateComment(r.Context(), req)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Created(w, response.M{
		"comment": comm,
	})
}

func (c *ArticleController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	iu := jwt.CurrentUser(r)
	if err := c.articleService.DeleteCommentByID(r.Context(), mux.Vars(r)["id"], iu); err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetArticleComments(w http.ResponseWriter, r *http.Request) {
	userID, err := jwt.GetUsernameFromReq(r)
	if err != nil {
		response.Err(w, err)
		return
	}

	comms, err := c.articleService.GetComments(r.Context(), mux.Vars(r)["slug"], userID)
	if err != nil {
		response.Err(w, err)
		return
	}

	response.Ok(w, response.M{
		"comments": comms,
	})
}
