package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/jmoiron/sqlx"
)

type ArticleRepoImpl struct {
	db *sqlx.DB
}

type ArticleRepository interface {
	InsertOne(context.Context, *dto.CreateArticleFields, string) (*model.Article, error)
	FindOneBySlug(context.Context, string) (*model.Article, error)
	DeleteBySlug(context.Context, string) error
	UpdateOneBySlug(context.Context, string, *UpdateArticleValues, *model.Article) error
	Find(context.Context, *FindArticlesArgs) (model.Articles, error)
	FindByFollowed(context.Context, *FindArticlesArgs) (model.Articles, error)
}

func (r *ArticleRepoImpl) InsertOne(ctx context.Context, d *dto.CreateArticleFields, authorID string) (*model.Article, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	a := &model.Article{
		Title:       d.Title,
		Description: d.Description,
		Body:        d.Body,
		AuthorID:    authorID,
	}

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
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.GetContext(ctx, a, *a); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *ArticleRepoImpl) FindOneBySlug(ctx context.Context, slug string) (*model.Article, error) {
	query := `
	SELECT
		id, title, description, body, author_id, created_at, updated_at
	FROM
		articles
	WHERE
		articles.slug = $1`
	a := new(model.Article)
	if err := r.db.GetContext(ctx, a, query, slug); err != nil {
		return nil, err
	}

	a.Slug = slug
	a.Author = new(model.User)
	return a, nil
}

func (r *ArticleRepoImpl) DeleteBySlug(ctx context.Context, slug string) error {
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

func (r *ArticleRepoImpl) UpdateOneBySlug(ctx context.Context, slug string, a *UpdateArticleValues, dest *model.Article) error {
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

func (r *ArticleRepoImpl) Find(ctx context.Context, p *FindArticlesArgs) (model.Articles, error) {
	articles := model.Articles{}
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

func (r *ArticleRepoImpl) FindByFollowed(ctx context.Context, p *FindArticlesArgs) (model.Articles, error) {
	articles := model.Articles{}
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
