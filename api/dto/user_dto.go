package dto

import "github.com/ashalfarhan/realworld/db/model"

type RegisterUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum,max=40"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginUserDto struct {
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Username string `json:"username" validate:"required_without=Email,omitempty,alphanum"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDto struct {
	Email    string           `json:"email" validate:"omitempty,email"`
	Username string           `json:"username" validate:"omitempty,alphanum,max=40"`
	Password string           `json:"password" validate:"omitempty,min=8,max=255"`
	Image    model.NullString `json:"image" validate:"url"`
	Bio      model.NullString `json:"bio" validate:"max=255"`
}
