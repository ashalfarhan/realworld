package repository

import "database/sql"

type FollowingRepository struct {
	db *sql.DB
}

func (r *FollowingRepository) Follow(follower, following string) error {
	_, err := r.db.Exec(`
	INSERT INTO 
		followings
		(follower_id, following_id)
	VALUES ($1, $2)
	`, follower, following)
	return err
}

func (r *FollowingRepository) Unfollow(follower, following string) error {
	_, err := r.db.Exec(`
	DELETE FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2
	`, follower, following)
	return err
}

func (r *FollowingRepository) IsFollowing(follower, following string) bool {
	return r.db.QueryRow(`
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
