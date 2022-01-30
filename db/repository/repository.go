package repository

import "database/sql"

type Repository struct {
	UserRepo        *UserRepository
	FollowRepo      *FollowingRepository
	ArticleRepo     *ArticleRepository
	ArticleTagsRepo *ArticleTagsRepository
}

func InitRepository(d *sql.DB) *Repository {
	return &Repository{
		UserRepo:        &UserRepository{d},
		FollowRepo:      &FollowingRepository{d},
		ArticleRepo:     &ArticleRepository{d},
		ArticleTagsRepo: &ArticleTagsRepository{d},
	}
}
