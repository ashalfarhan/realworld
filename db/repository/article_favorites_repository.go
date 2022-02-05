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

	query := `
	SELECT 
		article_id 
	FROM 
		article_favorites 
	WHERE 
		article_favorites.user_id = $1`

	if err := r.db.SelectContext(ctx, &article_ids, query, userID); err != nil {
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

	query := `
	INSERT INTO 
		article_favorites 
		(user_id, article_id)
	VALUES 
		($1, $2)`

	if _, err = tx.ExecContext(ctx, query, userID, articleID); err != nil {
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

	query := `
	DELETE FROM 
		article_favorites 
	WHERE 
		article_favorites.user_id = $1 
	AND 
		article_favorites.article_id = $2`

	if _, err = tx.ExecContext(ctx, query, userID, articleID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleFavoritesRepository) FindOneByIDs(ctx context.Context, userID, articleID string) (*string, error) {
	var ptr string
	query := `
	SELECT 
		article_favorites.user_id
	FROM 
		article_favorites 
	WHERE 
		article_favorites.user_id = $1 
	AND 
		article_favorites.article_id = $2`
	if err := r.db.QueryRowContext(ctx, query, userID, articleID).Scan(&ptr); err != nil {
		return nil, err
	}
	return &ptr, nil
}

func (r *ArticleFavoritesRepository) CountFavorites(ctx context.Context, articleID string) (int, error) {
	var count int
	query := `
	SELECT
		COUNT(*)
	FROM
		article_favorites
	WHERE
		article_favorites.article_id = $1
	`
	if err := r.db.QueryRowContext(ctx, query, articleID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
