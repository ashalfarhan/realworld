package model

type UserResponse struct {
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Bio      NullString `json:"bio"`
	Image    NullString `json:"image"`
	Token    string     `json:"token,omitempty"`
}

type ProfileResponse struct {
	Username  string     `json:"username"`
	Bio       NullString `json:"bio"`
	Image     NullString `json:"image"`
	Following bool       `json:"following"`
}

type ArticleResponse struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	TagList        []string `json:"tagList"`
	Author         *User    `json:"author"`
	Favorited      bool     `json:"favorited,omitempty"`
	FavoritesCount int      `json:"favoritesCount,omitempty"`
}
