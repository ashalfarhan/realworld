package dto

type CreateCommentFields struct {
	Body string `json:"body" validate:"required"`
}

type CreateCommentDto struct {
	AuthorID    string
	ArticleSlug string
	Comment     *CreateCommentFields `json:"comment" validate:"required"`
}
