package model

type CreateCommentFields struct {
	Body string `json:"body" validate:"required"`
}

type CreateCommentDto struct {
	Comment *CreateCommentFields `json:"comment" validate:"required"`
}
