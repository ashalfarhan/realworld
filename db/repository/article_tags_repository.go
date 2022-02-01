package repository

import (
	"context"
	"database/sql"
)

type ArticleTagsRepository struct {
	db *sql.DB
}

func (r *ArticleTagsRepository) InsertOne(ctx context.Context, articleID string, tagName string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, `
	INSERT INTO 
		article_tags
		(article_id, tag_name)
	VALUES
		($1, $2)
	`, articleID, tagName); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleTagsRepository) GetArticleTagsById(ctx context.Context, articleID string) ([]string, error) {
	row, err := r.db.QueryContext(ctx, `
	SELECT 
		tag_name
	FROM 
		article_tags
	WHERE
		article_tags.article_id = $1
	`, articleID)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var tags []string

	for row.Next() {
		var tag string
		if err := row.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *ArticleTagsRepository) GetAllTags(ctx context.Context) ([]string, error) {
	row, err := r.db.QueryContext(ctx, `
	SELECT 
		DISTINCT(tag_name)
	FROM 
		article_tags
	`)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	tags := []string{}

	for row.Next() {
		var tag string
		if err := row.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
