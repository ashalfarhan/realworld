package model

import "encoding/json"

type Comment struct {
	ID        string `json:"id" db:"id"`
	Body      string `json:"body" db:"body"`
	ArticleID string `json:"-" db:"article_id"`
	AuthorID  string `json:"-" db:"author_id"`
	Author    *User  `json:"author"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type Comments []*Comment

func (a Comments) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Comments) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
