package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleFavoritesRepoImpl struct {
	db *sqlx.DB
}

type ArticleFavoritesRepository interface {
	InsertOne(context.Context, string, string) error
	Delete(context.Context, string, string) error
	FindOneByIDs(context.Context, string, string) (*string, error)
	CountFavorites(context.Context, string) (int, error)
}

func (r *ArticleFavoritesRepoImpl) InsertOne(ctx context.Context, username, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO article_favorites (username, article_id) VALUES ($1, $2)"
	if _, err = tx.ExecContext(ctx, query, username, articleID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ArticleFavoritesRepoImpl) Delete(ctx context.Context, username, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	DELETE FROM article_favorites as af
	WHERE af.username = $1 
	AND af.article_id = $2`
	if _, err = tx.ExecContext(ctx, query, username, articleID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ArticleFavoritesRepoImpl) FindOneByIDs(ctx context.Context, username, articleID string) (*string, error) {
	var ptr string
	query := `
	SELECT af.username FROM article_favorites as af
	WHERE af.username = $1 
	AND af.article_id = $2`
	if err := r.db.QueryRowContext(ctx, query, username, articleID).Scan(&ptr); err != nil {
		return nil, err
	}
	return &ptr, nil
}

func (r *ArticleFavoritesRepoImpl) CountFavorites(ctx context.Context, articleID string) (int, error) {
	var count int
	query := `
	SELECT COUNT(*) FROM article_favorites as af
	WHERE af.article_id = $1`
	if err := r.db.QueryRowContext(ctx, query, articleID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
