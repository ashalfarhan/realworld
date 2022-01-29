package repository

import "database/sql"

type Repository struct {
	UR *UserRepository
	FR *FollowingRepository
}

func InitRepository(d *sql.DB) *Repository {
	return &Repository{
		UR: &UserRepository{d},
		FR: &FollowingRepository{d},
	}
}
