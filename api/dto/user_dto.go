package dto

import "github.com/ashalfarhan/realworld/db/model"

type RegisterUserFields struct {
	Email    string `json:"email" validate:"required,email" db:"email"`
	Username string `json:"username" validate:"required,max=40" db:"username"`
	Password string `json:"password" validate:"required,min=8,max=255" db:"password"`
}

type RegisterUserDto struct {
	User *RegisterUserFields `json:"user" validate:"required"`
}

type LoginUserFields struct {
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Username string `json:"username" validate:"required_without=Email,omitempty"`
	Password string `json:"password" validate:"required"`
}

type LoginUserDto struct {
	User *LoginUserFields `json:"user" validate:"required"`
}

type UpdateUserFields struct {
	Email    *string          `json:"email" validate:"omitempty,email"`
	Username *string          `json:"username" validate:"omitempty,max=40"`
	Password *string          `json:"password" validate:"omitempty,min=8,max=255"`
	Image    model.NullString `json:"image" validate:"url"`
	Bio      model.NullString `json:"bio" validate:"max=255"`
}

type UpdateUserDto struct {
	User *UpdateUserFields `json:"user" validate:"required"`
}
