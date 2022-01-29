package repository

import "database/sql"

type FollowingRepository struct {
	DB *sql.DB
}

func (f *FollowingRepository) Follow(follower, following string) error {
	_, err := f.DB.Exec(`
	INSERT INTO 
		followings
		(follower_id, following_id)
	VALUES ($1, $2)
	`, follower, following)
	return err
}

func (f *FollowingRepository) Unfollow(follower, following string) error {
	_, err := f.DB.Exec(`
	DELETE FROM
		followings
	WHERE
		followings.follower_id = $1
	AND
		followings.following_id = $2
	`, follower, following)
	return err
}

func (f *FollowingRepository) IsFollowing(follower, following string) bool {
	return f.DB.QueryRow(`
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
