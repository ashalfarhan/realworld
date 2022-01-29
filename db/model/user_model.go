package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Username  string     `json:"username"`
	Bio       NullString `json:"bio"`
	Image     NullString `json:"image"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (u *User) HashPassword(p string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) ValidatePassword(incPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(incPass)) == nil
}
