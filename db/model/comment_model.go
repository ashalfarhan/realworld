package model

type Comment struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author    User   `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
