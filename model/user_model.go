package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"-" db:"id"`
	Email     string     `json:"email" db:"email"`
	Password  string     `json:"-" db:"password"`
	Username  string     `json:"username" db:"username"`
	Bio       NullString `json:"bio" db:"bio"`
	Image     NullString `json:"image" db:"image"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"-" db:"updated_at"`
}

func (u *User) ValidatePassword(incPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(incPass)) == nil
}

type UserRs struct {
	Username string     `json:"username"`
	Bio      NullString `json:"bio"`
	Image    NullString `json:"image"`
	Email    string     `json:"email"`
	Token    string     `json:"token,omitempty"`
}

func (u *User) Serialize(token string) *UserRs {
	return &UserRs{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    token,
	}
}

type ProfileRs struct {
	Username  string     `json:"username"`
	Bio       NullString `json:"bio"`
	Image     NullString `json:"image"`
	Following bool       `json:"following"`
}

func (u *User) Profile(following bool) *ProfileRs {
	return &ProfileRs{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}
}

type FindUserArg struct {
	Email    string
	Username string
}
