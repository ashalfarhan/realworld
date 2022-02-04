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

	_, err = tx.ExecContext(ctx, `
	INSERT INTO
		article_tags
		(article_id, tag_name)
	VALUES
		($1, $2)`, articleID, tagName)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleTagsRepository) GetArticleTagsByID(ctx context.Context, articleID string) ([]string, error) {
	var tags []string

	err := r.db.SelectContext(ctx, &tags, `
	SELECT
		tag_name
	FROM
		article_tags
	WHERE
		article_tags.article_id = $1`, articleID)

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *ArticleTagsRepository) GetAllTags(ctx context.Context) ([]string, error) {
	var tags []string

	err := r.db.SelectContext(ctx, &tags, `
	SELECT
		DISTINCT(tag_name)
	FROM
		article_tags`)

	if err != nil {
		return nil, err
	}

	return tags, nil
}
