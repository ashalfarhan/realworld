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
