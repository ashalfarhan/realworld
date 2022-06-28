package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleTagsRepo struct {
	db *sqlx.DB
}

type ArticleTagsRepository interface {
	InsertBulk(ctx context.Context, tags []InsertArticleTagsArgs) error
	FindArticleTagsByID(ctx context.Context, articleID string) ([]string, error)
	FindAllTags(ctx context.Context) ([]string, error)
}

type InsertArticleTagsArgs struct {
	ArticleID string `db:"article_id"`
	TagName   string `db:"tag_name"`
}

func (r *ArticleTagsRepo) InsertBulk(ctx context.Context, tags []InsertArticleTagsArgs) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO article_tags (article_id, tag_name)
	VALUES (:article_id, :tag_name)`
	if _, err = tx.NamedExecContext(ctx, query, tags); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ArticleTagsRepo) FindArticleTagsByID(ctx context.Context, articleID string) ([]string, error) {
	var tags []string

	query := `
	SELECT tag_name FROM article_tags as at
	WHERE at.article_id = $1
	ORDER BY tag_name ASC`
	if err := r.db.SelectContext(ctx, &tags, query, articleID); err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *ArticleTagsRepo) FindAllTags(ctx context.Context) ([]string, error) {
	var tags []string
	query := "SELECT DISTINCT(tag_name) FROM article_tags ORDER BY tag_name ASC"
	if err := r.db.SelectContext(ctx, &tags, query); err != nil {
		return nil, err
	}
	return tags, nil
}
