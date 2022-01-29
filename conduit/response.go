package conduit

import "github.com/ashalfarhan/realworld/db/model"

type UserAuthResponse struct {
	Email    string           `json:"email"`
	Token    string           `json:"token"`
	Username string           `json:"username"`
	Bio      model.NullString `json:"bio"`
	Image    model.NullString `json:"image"`
}

type ProfileResponse struct {
	Username  string           `json:"username"`
	Bio       model.NullString `json:"bio"`
	Image     model.NullString `json:"Image"`
	Following bool             `json:"following"`
}
