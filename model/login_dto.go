package model

type LoginUserFields struct {
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Username string `json:"username" validate:"required_without=Email,omitempty"`
	Password string `json:"password" validate:"required,max=64"`
}

type LoginUserDto struct {
	User *LoginUserFields `json:"user" validate:"required"`
}
