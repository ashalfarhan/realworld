package dto

type CreateArticleDto struct {
	Title       string   `json:"title" validate:"required,max=255"`
	Description string   `json:"description" validate:"required,max=255"`
	Body        string   `json:"body" validate:"required"`
	TagList     []string `json:"tagList,omitempty" validate:"omitempty,unique"`
}
