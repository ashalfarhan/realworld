package model

type RegisterUserFields struct {
	Email    string `json:"email" validate:"required,email" db:"email"`
	Username string `json:"username" validate:"required,max=40" db:"username"`
	Password string `json:"password" validate:"required,min=8,max=64" db:"password"`
}

type RegisterUserDto struct {
	User *RegisterUserFields `json:"user" validate:"required"`
}
