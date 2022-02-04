package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	UserRepo             *UserRepository
	FollowRepo           *FollowingRepository
	ArticleRepo          *ArticleRepository
	ArticleTagsRepo      *ArticleTagsRepository
	ArticleFavoritesRepo *ArticleFavoritesRepository
}

func InitRepository(d *sqlx.DB) *Repository {
	return &Repository{
		UserRepo:             &UserRepository{d},
		FollowRepo:           &FollowingRepository{d},
		ArticleRepo:          &ArticleRepository{d},
		ArticleTagsRepo:      &ArticleTagsRepository{d},
		ArticleFavoritesRepo: &ArticleFavoritesRepository{d},
	}
}
