package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/model"
	"github.com/jmoiron/sqlx"
)

type ArticleRepoImpl struct {
	db *sqlx.DB
}

type ArticleRepository interface {
	InsertOne(context.Context, *model.CreateArticleFields, string) (*model.Article, error)
	FindOneBySlug(context.Context, string, string) (*model.Article, error)
	DeleteBySlug(context.Context, string) error
	UpdateOneBySlug(context.Context, *model.UpdateArticleFields, *model.Article) error
	Find(context.Context, *model.FindArticlesArgs) (model.Articles, error)
}

func (r *ArticleRepoImpl) InsertOne(ctx context.Context, d *model.CreateArticleFields, username string) (*model.Article, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	a := &model.Article{
		Title:          d.Title,
		Description:    d.Description,
		Body:           d.Body,
		AuthorUsername: username,
		Slug:           d.Slug,
	}

	query := `
	INSERT INTO articles (slug, title, description, body, author_username) 
	VALUES (:slug, :title, :description, :body, :author_username) 
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

func (r *ArticleRepoImpl) FindOneBySlug(ctx context.Context, username, slug string) (*model.Article, error) {
	query := `
	SELECT
		ar.id, ar.author_username, ar.title, ar.description, ar.body, 
		ar.created_at, ar.updated_at, ar.slug,
		us.username as "author.username", us.bio as "author.bio",
		us.image as "author.image", COUNT(af.username) as "favorites_count",
		($1 IN (af.username) IS NOT NULL) as "favorited"
	FROM articles as ar 
	LEFT JOIN users as us
		ON us.username = ar.author_username
	LEFT JOIN article_favorites as af
		ON af.article_id = ar.id
	WHERE ar.slug = $2
	GROUP BY (ar.id, us.username, us.bio, us.image, af.username)`
	a := new(model.Article)
	if err := r.db.GetContext(ctx, a, query, username, slug); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *ArticleRepoImpl) Find(ctx context.Context, p *model.FindArticlesArgs) (model.Articles, error) {
	articles := model.Articles{}
	query := `
	SELECT 
		ar.id, ar.author_username, ar.title, ar.description, ar.body, 
		ar.created_at, ar.updated_at, ar.slug,
		us.username as "author.username", us.bio as "author.bio",
		us.image as "author.image", 
		COUNT(af.username) as "favorites_count", 
		(:username IN (af.username) IS NOT NULL) as "favorited"
	FROM articles as ar
	LEFT JOIN users as us
		ON us.username = ar.author_username
	LEFT JOIN article_favorites as af
		ON af.article_id = ar.id
	WHERE 1 = 1`

	if p.Author != "" {
		query += `
		AND ar.author_username = :author_username`
	}

	if p.Tag != "" {
		query += `
		AND ar.id IN (
			SELECT at.article_id 
			FROM article_tags as at
			WHERE at.tag_name = :tag
		)`
	}

	if p.Favorited != "" {
		query += `
		AND ar.id IN (
			SELECT af.article_id
			FROM article_favorites as af
			WHERE af.username = :favorited_by
		)`
	}

	if p.Username != "" {
		query += `
		AND ar.author_username IN (
			SELECT f.following_username 
			FROM followings as f
			WHERE f.follower_username = :username
		)`
	}

	query += " GROUP BY (ar.id, us.username, us.bio, us.image, af.username) ORDER BY ar.created_at DESC LIMIT :limit OFFSET :offset"
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
