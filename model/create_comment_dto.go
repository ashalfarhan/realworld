package model

type CreateCommentFields struct {
	Body string `json:"body" validate:"required"`
}

type CreateCommentDto struct {
	AuthorUsername,
	ArticleSlug string
	Comment *CreateCommentFields `json:"comment" validate:"required"`
}
