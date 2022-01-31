package dto

import "github.com/ashalfarhan/realworld/db/model"

type RegisterUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,alphanum,max=40"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginUserDto struct {
	Email    string `json:"email,omitempty" validate:"required_without=Username,omitempty,email"`
	Username string `json:"username,omitempty" validate:"required_without=Email,omitempty,alphanum"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDto struct {
	Email    string           `json:"email,omitempty" validate:"omitempty,email"`
	Username string           `json:"username,omitempty" validate:"omitempty,alphanum,max=40"`
	Password string           `json:"password,omitempty" validate:"omitempty,min=8,max=255"`
	Image    model.NullString `json:"image,omitempty" validate:"omitempty,max=255,url"`
	Bio      model.NullString `json:"bio,omitempty" validate:"omitempty,max=255"`
}
