package dto

type CreateCommentFields struct {
	Body        string `json:"body" validate:"required"`
	AuthorID    string
	ArticleSlug string
}

type CreateCommentDto struct {
	Comment *CreateCommentFields `json:"comment" validate:"required"`
}
