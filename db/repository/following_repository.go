package repository

import (
	"context"
	"database/sql"
)

type FollowingRepository struct {
	db *sql.DB
}

func (r *FollowingRepository) Follow(ctx context.Context, follower, following string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	INSERT INTO 
		followings
		(follower_id, following_id)
	VALUES ($1, $2)
	`, follower, following)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FollowingRepository) Unfollow(ctx context.Context, follower, following string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	DELETE FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2
	`, follower, following)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FollowingRepository) IsFollowing(ctx context.Context, follower, following string) bool {
	return r.db.QueryRowContext(ctx, `
	SELECT 
		COUNT(*)
	FROM 
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2
	`, follower, following).Err() == nil

}
