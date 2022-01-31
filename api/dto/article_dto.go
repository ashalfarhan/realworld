package dto

type CreateArticleDto struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required,max=255"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList" validate:"omitempty,unique"`
}

type UpdateArticleDto struct {
	Title       string `json:"title" validate:"max=255"`
	Body        string `json:"body"`
	Description string `json:"description" validate:"max=255"`
}
