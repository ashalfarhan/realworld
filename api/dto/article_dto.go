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
	Title       string `json:"title" validate:"max=255"`
	Body        string `json:"body"`
	Description string `json:"description" validate:"max=255"`
}

type UpdateArticleDto struct {
	Article *UpdateArticleFields `json:"article" validate:"required"`
}
