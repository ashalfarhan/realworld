package model

type UpdateUserFields struct {
	Email    *string    `json:"email" validate:"omitempty,email"`
	Username *string    `json:"username" validate:"omitempty,max=40"`
	Password *string    `json:"password" validate:"omitempty,min=8,max=64"`
	Image    NullString `json:"image" validate:"url"`
	Bio      NullString `json:"bio" validate:"max=255"`
}

type UpdateUserDto struct {
	User *UpdateUserFields `json:"user" validate:"required"`
}
