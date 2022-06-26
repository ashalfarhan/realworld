package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/model"
	"github.com/ashalfarhan/realworld/utils/logger"
	"github.com/jmoiron/sqlx"
)

type ArticleRepoImpl struct {
	db *sqlx.DB
}

type ArticleRepository interface {
	InsertOne(context.Context, *model.CreateArticleFields, string) (*model.Article, error)
	FindOneBySlug(context.Context, string) (*model.Article, error)
	DeleteBySlug(context.Context, string) error
	UpdateOneBySlug(context.Context, *model.UpdateArticleFields, *model.Article) error
	Find(context.Context, *FindArticlesArgs) (model.Articles, error)
}

func (r *ArticleRepoImpl) InsertOne(ctx context.Context, d *model.CreateArticleFields, authorID string) (*model.Article, error) {
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
		Slug:        d.Slug,
	}

	query := `
	INSERT INTO articles (slug, title, description, body, author_id) 
	VALUES (:slug, :title, :description, :body, :author_id) 
	RETURNING id, created_at, updated_at`
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

func (r *ArticleRepoImpl) DeleteBySlug(ctx context.Context, slug string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "DELETE FROM articles as a WHERE a.slug = $1"

	if _, err = tx.ExecContext(ctx, query, slug); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepoImpl) UpdateOneBySlug(ctx context.Context, d *model.UpdateArticleFields, a *model.Article) error {
	if v := d.Title; v != nil {
		a.Title = *v
	}
	if v := d.Slug; v != nil {
		a.Slug = *v
	}
	if v := d.Body; v != nil {
		a.Body = *v
	}
	if v := d.Description; v != nil {
		a.Description = *v
	}

	query := `
	UPDATE articles as a
	SET 
		title = :title, slug = :slug,
		body = :body, description = :description,
		updated_at = NOW()
	WHERE a.id = :id`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, a); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleRepoImpl) FindOneBySlug(ctx context.Context, slug string) (*model.Article, error) {
	query := `
	SELECT id, title, description, body, author_id, created_at, updated_at, slug
	FROM articles as a
	WHERE a.slug = $1`
	a := new(model.Article)
	if err := r.db.GetContext(ctx, a, query, slug); err != nil {
		return nil, err
	}

	a.Author = new(model.User)
	return a, nil
}

type FindArticlesArgs struct {
	Tag       string `db:"tag"`
	Author    string `validate:"omitempty,uuid" db:"author_id"` // TODO: Change to username
	UserID    string `db:"user_id"`
	Favorited string `db:"favorited_by"` // TODO: Change to username
	Limit     int    `validate:"min=1,max=25" db:"limit"`
	Offset    int    `validate:"min=0" db:"offset"`
}

func (r *ArticleRepoImpl) Find(ctx context.Context, p *FindArticlesArgs) (model.Articles, error) {
	articles := model.Articles{}
	query := `
	SELECT 
		a.id, a.title, 
		a.description, a.body, 
		a.author_id, a.created_at, 
		a.updated_at, a.slug 
	FROM articles as a WHERE 1 = 1`

	if p.Author != "" {
		query += `
		AND a.author_id = :author_id`
	}

	if p.Tag != "" {
		query += `
		AND a.id IN (
			SELECT at.article_id 
			FROM article_tags as at
			WHERE at.tag_name = :tag
		)`
	}

	if p.Favorited != "" {
		query += `
		AND a.id IN (
			SELECT af.article_id
			FROM article_favorites as af
			WHERE af.user_id = :favorited_by
		)`
	}

	if p.UserID != "" {
		query += `
		AND a.author_id IN (
			SELECT f.following_id 
			FROM followings as f
			WHERE f.follower_id = :user_id
		)`
	}

	query += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	logger.Log.Printf("Find sql:%s", query)
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
