package repository

import (
	"database/sql"
	"fmt"
	"strings"

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

func (r *ArticleRepository) Update(slug string, a *conduit.UpdateArticleArgs) error {
	var updateArgs []string
	var valArgs []interface{}
	argIdx := 0

	if v := a.Body; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("body = $%d", argIdx))
		valArgs = append(valArgs, *a.Body)
	}

	if v := a.Title; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("title = $%d", argIdx))
		valArgs = append(valArgs, *a.Title)
	}

	if v := a.Slug; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("slug = $%d", argIdx))
		valArgs = append(valArgs, *a.Slug)
	}

	if v := a.Description; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("description = $%d", argIdx))
		valArgs = append(valArgs, *a.Description)
	}

	updateArgs = append(updateArgs, "updated_at = NOW()")

	argIdx++
	valArgs = append(valArgs, slug)
	query := fmt.Sprintf("UPDATE articles SET %s WHERE articles.slug = $%d", strings.Join(updateArgs, ", "), argIdx)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if _, err := stmt.Exec(valArgs...); err != nil {
		return err
	}

	return nil
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
		if err = row.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.Body,
			&a.Author.ID,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}
