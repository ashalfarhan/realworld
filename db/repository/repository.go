package repository

import "database/sql"

type Repository struct {
	UserRepo        *UserRepository
	FollowRepo      *FollowingRepository
	ArticleRepo     *ArticleRepository
	TagRepo         *TagRepository
	ArticleTagsRepo *ArticleTagsRepository
}

func InitRepository(d *sql.DB) *Repository {
	return &Repository{
		UserRepo:        &UserRepository{d},
		FollowRepo:      &FollowingRepository{d},
		ArticleRepo:     &ArticleRepository{d},
		TagRepo:         &TagRepository{d},
		ArticleTagsRepo: &ArticleTagsRepository{d},
	}
}
