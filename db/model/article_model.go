package model

type Article struct {
	ID          string   `json:"-" db:"id"`
	Slug        string   `json:"slug" db:"slug"`
	Title       string   `json:"title" db:"title"`
	Description string   `json:"description" db:"description"`
	Body        string   `json:"body" db:"body"`
	CreatedAt   string   `json:"createdAt" db:"created_at"`
	UpdatedAt   string   `json:"updatedAt" db:"updated_at"`
	TagList     []string `json:"tagList"`
	AuthorID    string   `json:"-" db:"author_id"`
	Author      *User    `json:"author" db:"author"`
}
