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
	err := utils.ValidateDTO(r, req)
	if err != nil {
		response.Err(w, err)
		return
	}
	iu, _ := jwt.CurrentUser(r)

	req.AuthorID = iu.Subject
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
	iu, _ := jwt.CurrentUser(r)
	if err := c.articleService.DeleteCommentByID(r.Context(), mux.Vars(r)["id"], iu.Subject); err != nil {
		response.Err(w, err)
		return
	}

	response.Accepted(w, nil)
}

func (c *ArticleController) GetArticleComments(w http.ResponseWriter, r *http.Request) {
	userID, err := jwt.GetUserIDFromReq(r)
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
