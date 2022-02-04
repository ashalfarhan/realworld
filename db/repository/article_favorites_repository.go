package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleFavoritesRepository struct {
	db *sqlx.DB
}

func (r *ArticleFavoritesRepository) FindFavoritedArticleByUserId(ctx context.Context, userID string) ([]string, error) {
	var article_ids []string

	err := r.db.SelectContext(ctx, &article_ids, `
	SELECT 
		article_id 
	FROM 
		article_favorites 
	WHERE 
		article_favorites.user_id = $1`, userID)

	if err != nil {
		return nil, err
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
		($1, $2)`, userID, articleID)

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
		article_favorites.article_id = $2`, userID, articleID)

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
		article_favorites.article_id = $2`, userID, articleID).Err()
}
