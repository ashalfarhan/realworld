package conduit

import "github.com/ashalfarhan/realworld/db/model"

type ArticleParams struct {
	Tag    string `validate:"max=20"`
	Author string `validate:"alphanum"`
	Limit  int    `validate:"min=1,max=25"`
	Offset int    `validate:"min=0"`
}

// Use Pointer to update
// To determine if the field needs to be updated
type UpdateUserArgs struct {
	ID       string
	Email    *string
	Username *string
	Password *string
	Image    model.NullString
	Bio      model.NullString
}

type UpdateArticleArgs struct {
	Title       *string
	Slug        *string
	Body        *string
	Description *string
}
