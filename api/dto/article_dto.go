package dto

type CreateArticleFields struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required,max=255"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList" validate:"omitempty,unique"`
	Slug        string
}

type CreateArticleDto struct {
	Article *CreateArticleFields `json:"article" validate:"required"`
}

type UpdateArticleFields struct {
	Title       *string `json:"title" validate:"omitempty,max=255"`
	Body        *string `json:"body" validate:"omitempty,max=255"`
	Description *string `json:"description" validate:"omitempty,max=255"`
	Slug        *string
}

type UpdateArticleDto struct {
	Article *UpdateArticleFields `json:"article" validate:"required"`
}
