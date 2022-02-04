package conduit

import "github.com/ashalfarhan/realworld/db/model"

type UserResponse struct {
	Email    string           `json:"email"`
	Username string           `json:"username"`
	Bio      model.NullString `json:"bio"`
	Image    model.NullString `json:"image"`
	Token    string           `json:"token,omitempty"`
}

type ProfileResponse struct {
	Username  string           `json:"username"`
	Bio       model.NullString `json:"bio"`
	Image     model.NullString `json:"image"`
	Following bool             `json:"following"`
}

type ArticleResponse struct {
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	CreatedAt      string      `json:"createdAt"`
	UpdatedAt      string      `json:"updatedAt"`
	TagList        []string    `json:"tagList"`
	Author         *model.User `json:"author"`
	Favorited      bool        `json:"favorited,omitempty"`
	FavoritesCount int         `json:"favoritesCount,omitempty"`
}
