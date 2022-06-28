package model

import (
	"encoding/json"
	"time"
)

type Article struct {
	ID             string     `json:"id" db:"id"`
	Slug           string     `json:"slug" db:"slug"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Body           string     `json:"body" db:"body"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
	TagList        []string   `json:"tagList"`
	AuthorUsername string     `json:"authorUsername" db:"author_username"`
	Favorited      bool       `json:"favorited" db:"favorited"`
	FavoritesCount int        `json:"favoritesCount" db:"favorites_count"`
	Author         *ProfileRs `json:"author" db:"author"`
}

type ArticleRs struct {
	Slug           string     `json:"slug"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Body           string     `json:"body"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	TagList        []string   `json:"tagList"`
	Favorited      bool       `json:"favorited"`
	FavoritesCount int        `json:"favoritesCount"`
	Author         *ProfileRs `json:"author"`
}

func (a Article) Serialize() *ArticleRs {
	return &ArticleRs{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		TagList:        a.TagList,
		Favorited:      a.Favorited,
		FavoritesCount: a.FavoritesCount,
		Author:         a.Author,
	}
}

func (a Article) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Article) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

type Articles []*Article

func (as Articles) Serialize() []*ArticleRs {
	ars := []*ArticleRs{}
	for _, a := range as {
		ars = append(ars, a.Serialize())
	}
	return ars
}

func (a Articles) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a *Articles) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
