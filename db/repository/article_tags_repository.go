package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleTagsRepository struct {
	db *sqlx.DB
}

func (r *ArticleTagsRepository) InsertOne(ctx context.Context, articleID, tagName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO
		article_tags
		(article_id, tag_name)
	VALUES
		($1, $2)`

	if _, err = tx.ExecContext(ctx, query, articleID, tagName); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleTagsRepository) FindArticleTagsByID(ctx context.Context, articleID string) ([]string, error) {
	var tags []string

	query := `
	SELECT
		tag_name
	FROM
		article_tags
	WHERE
		article_tags.article_id = $1
	ORDER BY tag_name ASC`
	if err := r.db.SelectContext(ctx, &tags, query, articleID); err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *ArticleTagsRepository) FindAllTags(ctx context.Context) ([]string, error) {
	var tags []string

	query := `
	SELECT
		DISTINCT(tag_name)
	FROM
		article_tags
	ORDER BY tag_name ASC`
	if err := r.db.SelectContext(ctx, &tags, query); err != nil {
		return nil, err
	}

	return tags, nil
}
