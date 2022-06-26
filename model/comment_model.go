package model

import (
	"encoding/json"
	"time"
)

type Comment struct {
	ID        string    `json:"id" db:"id"`
	Body      string    `json:"body" db:"body"`
	ArticleID string    `json:"-" db:"article_id"`
	AuthorID  string    `json:"-" db:"author_id"`
	Author    *User     `json:"author"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Comments []*Comment

func (a Comments) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Comments) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
