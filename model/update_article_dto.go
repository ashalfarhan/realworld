package model

type UpdateArticleFields struct {
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Body        *string `json:"body" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=255"`
	Slug        *string
}

type UpdateArticleDto struct {
	Article *UpdateArticleFields `json:"article" validate:"required"`
}
