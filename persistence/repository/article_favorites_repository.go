package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ArticleFavoritesRepoImpl struct {
	db *sqlx.DB
}

type ArticleFavoritesRepository interface {
	FindFavoritedArticleByUserId(context.Context, string) ([]string, error)
	InsertOne(context.Context, string, string) error
	Delete(context.Context, string, string) error
	FindOneByIDs(context.Context, string, string) (*string, error)
	CountFavorites(context.Context, string) (int, error)
}

func (r *ArticleFavoritesRepoImpl) FindFavoritedArticleByUserId(ctx context.Context, userID string) ([]string, error) {
	var article_ids []string

	query := `
	SELECT article_id FROM article_favorites 
	WHERE article_favorites.user_id = $1`
	if err := r.db.SelectContext(ctx, &article_ids, query, userID); err != nil {
		return nil, err
	}

	return article_ids, nil
}

func (r *ArticleFavoritesRepoImpl) InsertOne(ctx context.Context, userID, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO article_favorites (user_id, article_id) VALUES ($1, $2)"
	if _, err = tx.ExecContext(ctx, query, userID, articleID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleFavoritesRepoImpl) Delete(ctx context.Context, userID, articleID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	DELETE FROM article_favorites as af
	WHERE af.user_id = $1 
	AND af.article_id = $2`
	if _, err = tx.ExecContext(ctx, query, userID, articleID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ArticleFavoritesRepoImpl) FindOneByIDs(ctx context.Context, userID, articleID string) (*string, error) {
	var ptr string
	query := `
	SELECT af.user_id FROM article_favorites as af
	WHERE af.user_id = $1 
	AND af.article_id = $2`
	if err := r.db.QueryRowContext(ctx, query, userID, articleID).Scan(&ptr); err != nil {
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
