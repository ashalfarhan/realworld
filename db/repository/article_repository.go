package repository

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/db/model"
)

type ArticleRepository struct {
	db *sql.DB
}

func (r *ArticleRepository) InsertOne(a *model.Article) error {
	return r.db.QueryRow(`
	INSERT INTO 
		articles
		(slug, title, description, body, author_id)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		articles.created_at,
		articles.updated_at
	`).Scan(&a.CreatedAt, &a.UpdatedAt)
}

func (r *ArticleRepository) FindOneBySlug(a *model.Article) error {
	return r.db.QueryRow(`
	SELECT 
		title, description, body, author_id
	FROM 
		articles
	WHERE 
		articles.slug = $1
	`, a.Slug).Scan(&a.Title, &a.Description, &a.Body, &a.Author.ID)
}

func (r *ArticleRepository) DeleteBySlug(slug string) error {
	_, err := r.db.Exec(`
	DELETE FROM
		articles
	WHERE 
		articles.slug = $1
	`, slug)

	return err
}
