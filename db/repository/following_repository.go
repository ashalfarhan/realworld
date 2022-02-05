package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type FollowingRepository struct {
	db *sqlx.DB
}

func (r *FollowingRepository) InsertOne(ctx context.Context, follower, following string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
	INSERT INTO 
		followings
		(follower_id, following_id)
	VALUES ($1, $2)`
	if _, err = tx.ExecContext(ctx, query, follower, following); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FollowingRepository) DeleteOneIDs(ctx context.Context, follower, following string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `
	DELETE FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2`
	if _, err = tx.ExecContext(ctx, query, follower, following); err != nil {
		return err
	}

	return tx.Commit()
}

// Returns pointer to the following id.
// To determine if "follower" is follow "following".
// Check if pointer is not nill and err is nil
func (r *FollowingRepository) FindOneByIDs(ctx context.Context, follower, following string) (*string, error) {
	var ptr *string
	query := `
	SELECT 
		followings.following_id
	FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2`
	if err := r.db.QueryRowContext(ctx, query, follower, following).Scan(ptr); err != nil {
		return nil, err
	}

	return ptr, nil
}
