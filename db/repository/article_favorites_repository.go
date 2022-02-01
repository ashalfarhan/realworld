package repository

import (
	"context"
	"database/sql"
)

type ArticleFavoritesRepository struct {
	db *sql.DB
}

func (r *ArticleFavoritesRepository) FindFavoritedArticleByUserId(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
	SELECT
		article_id
	FROM
		article_favorites
	WHERE
		article_favorites.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	article_ids := []string{}

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		article_ids = append(article_ids, id)
	}

	return article_ids, nil
}

func (r *ArticleFavoritesRepository) InsertOne(ctx context.Context, userID, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	INSERT INTO
		article_favorites
		(user_id, article_id)
	VALUES
		($1, $2)
	`, userID, articleID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleFavoritesRepository) Delete(ctx context.Context, userID, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	DELETE FROM
		article_favorites
	WHERE
		article_favorites.user_id = $1
	AND
		article_favorites.article_id = $2
	`, userID, articleID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleFavoritesRepository) GetOneByIDs(ctx context.Context, userID, articleID string) error {
	return r.db.QueryRowContext(ctx, `
	SELECT *
	FROM
		article_favorites
	WHERE
		article_favorites.user_id = $1
	AND
		article_favorites.article_id = $2
	`, userID, articleID).Err()
}
