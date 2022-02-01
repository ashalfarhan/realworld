package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
)

type ArticleRepository struct {
	db *sql.DB
}

func (r *ArticleRepository) InsertOne(ctx context.Context, a *model.Article) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
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

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) FindOneBySlug(ctx context.Context, a *model.Article) error {
	return r.db.QueryRowContext(ctx, `
	SELECT
		id, title, description, body, author_id, created_at, updated_at
	FROM
		articles
	WHERE
		articles.slug = $1
	`, a.Slug).Scan(&a.ID, &a.Title, &a.Description, &a.Body, &a.Author.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *ArticleRepository) DeleteBySlug(ctx context.Context, slug string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	DELETE FROM
		articles
	WHERE
		articles.slug = $1
	`, slug)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) UpdateOneBySlug(ctx context.Context, slug string, a *conduit.UpdateArticleArgs, dest *model.Article) error {
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
	query := fmt.Sprintf("UPDATE articles SET %s WHERE articles.slug = $%d RETURNING articles.updated_at", strings.Join(updateArgs, ", "), argIdx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, valArgs...).Scan(&dest.UpdatedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) Find(ctx context.Context, p *conduit.ArticleArgs) ([]*model.Article, error) {
	row, err := r.db.QueryContext(ctx, `
	SELECT
		id, title, description, body, author_id, created_at, updated_at, slug
	FROM
		articles
	ORDER BY created_at ASC
	LIMIT $1
	OFFSET $2
	`, p.Limit, p.Offset)

	if err != nil {
		return nil, err
	}

	articles := []*model.Article{}
	defer row.Close()

	for row.Next() {
		a := &model.Article{
			Author: &model.User{},
		}

		err = row.Scan(&a.ID, &a.Title, &a.Description, &a.Body, &a.Author.ID, &a.CreatedAt, &a.UpdatedAt, &a.Slug)

		if err != nil {
			return nil, err
		}

		articles = append(articles, a)
	}

	return articles, nil
}
