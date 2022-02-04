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

	_, err = tx.ExecContext(ctx, `
	INSERT INTO 
		followings
		(follower_id, following_id)
	VALUES ($1, $2)`, follower, following)

	if err != nil {
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

	_, err = tx.ExecContext(ctx, `
	DELETE FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2`, follower, following)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *FollowingRepository) GetOneByIDs(ctx context.Context, follower, following string) error {
	return r.db.QueryRowContext(ctx, `
	SELECT *
	FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2`, follower, following).Err()
}
