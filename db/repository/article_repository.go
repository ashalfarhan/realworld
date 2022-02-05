package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/jmoiron/sqlx"
)

type ArticleRepository struct {
	db *sqlx.DB
}

func (r *ArticleRepository) InsertOne(ctx context.Context, a *model.Article) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO 
		articles 
		(slug, title, description, body, author_id) 
	VALUES 
		(:slug, :title, :description, :body, :author_id) 
	RETURNING 
		id, created_at, updated_at`
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err = stmt.GetContext(ctx, a, *a); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepository) FindOneBySlug(ctx context.Context, a *model.Article) error {
	query := `
	SELECT
		id, title, description, body, author_id, created_at, updated_at
	FROM
		articles
	WHERE
		articles.slug = $1`

	return r.db.GetContext(ctx, a, query, a.Slug)
}

func (r *ArticleRepository) DeleteBySlug(ctx context.Context, slug string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	DELETE FROM 
		articles 
	WHERE 
		articles.slug = $1`

	if _, err = tx.ExecContext(ctx, query, slug); err != nil {
		return err
	}

	return tx.Commit()
}

type UpdateArticleValues struct {
	Title       *string
	Slug        *string
	Body        *string
	Description *string
}

func (r *ArticleRepository) UpdateOneBySlug(ctx context.Context, slug string, a *UpdateArticleValues, dest *model.Article) error {
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

type FindArticlesArgs struct {
	Tag string `db:"tag"`
	// Author string `validate:"alphanum" db:"author_id"`
	UserID string `db:"user_id"`
	Limit  int    `validate:"min=1,max=25" db:"limit"`
	Offset int    `validate:"min=0" db:"offset"`
}

func (r *ArticleRepository) Find(ctx context.Context, p *FindArticlesArgs) ([]*model.Article, error) {
	articles := []*model.Article{}
	query := "SELECT articles.id, articles.title, articles.description, articles.body, articles.author_id, articles.created_at, articles.updated_at, articles.slug FROM articles"
	if p.Tag != "" {
		query += ` 
		WHERE
			articles.id
		IN (
			SELECT 
				article_tags.article_id 
			FROM 
				article_tags
			WHERE
				article_tags.tag_name = :tag
			)`
	}

	query += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &articles, p); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *ArticleRepository) FindByFollowed(ctx context.Context, p *FindArticlesArgs) ([]*model.Article, error) {
	articles := []*model.Article{}
	query := "SELECT articles.id, articles.title, articles.description, articles.body, articles.author_id, articles.created_at, articles.updated_at, articles.slug FROM articles"

	if p.UserID != "" {
		query += ` 
		WHERE
			articles.author_id
		IN (
			SELECT 
				followings.following_id 
			FROM 
				followings
			WHERE
				followings.follower_id = :user_id
			)`
	}

	query += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	if err := stmt.SelectContext(ctx, &articles, p); err != nil {
		return nil, err
	}

	return articles, nil
}
