package model

type Article struct {
	ID          string   `json:"-"`
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	TagList     []string `json:"tagList"`
	Author      *User    `json:"author"`
}
