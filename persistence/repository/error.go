package repository

const (
	ErrDuplicateEmail     = "pq: duplicate key value violates unique constraint \"users_email_key\""
	ErrDuplicateUsername  = "pq: duplicate key value violates unique constraint \"users_username_key\""
	ErrDuplicateFollowing = "pq: duplicate key value violates unique constraint \"followings_pkey\""
)
