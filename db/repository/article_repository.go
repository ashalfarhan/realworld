package repository

import (
	"database/sql"

	"github.com/ashalfarhan/realworld/conduit"
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
		articles.id,
		articles.created_at,
		articles.updated_at
	`, a.Slug, a.Title, a.Description, a.Body, a.Author.ID).
		Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *ArticleRepository) FindOneBySlug(a *model.Article) error {
	return r.db.QueryRow(`
	SELECT
		id, title, description, body, author_id, created_at, updated_at
	FROM
		articles
	WHERE
		articles.slug = $1
	`, a.Slug).Scan(&a.ID, &a.Title, &a.Description, &a.Body, &a.Author.ID, &a.CreatedAt, &a.UpdatedAt)
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

func (r *ArticleRepository) GetFiltered(p *conduit.ArticleParams) ([]model.Article, error) {
	row, err := r.db.Query(`
	SELECT
		id, title, description, body, author_id, created_at, updated_at
	FROM
		articles
	LIMIT $1
	OFFSET $2
	ORDER BY created_at DESC
	`, p.Limit, p.Offset)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	defer row.Close()

	for row.Next() {
		a := model.Article{
			Author: &model.User{},
		}
		err = row.Scan(&a.ID, &a.Title, &a.Description, &a.Body, &a.Author.ID, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}
