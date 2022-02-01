package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"-"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Username  string     `json:"username"`
	Bio       NullString `json:"bio"`
	Image     NullString `json:"image"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}

func (u *User) ValidatePassword(incPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(incPass)) == nil
}
