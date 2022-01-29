package model

type Article struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	// Favorited      bool     `json:"favorited"`
	FavoritesCount int  `json:"favoritesCount"`
	Author         User `json:"author"`
}
