package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	UserRepo             UserRepository
	FollowRepo           FollowingRepository
	ArticleRepo          ArticleRepository
	ArticleTagsRepo      *ArticleTagsRepository
	ArticleFavoritesRepo ArticleFavoritesRepository
	CommentRepo          CommentRepository
}

func InitRepository(d *sqlx.DB) *Repository {
	return &Repository{
		&UserRepoImpl{d},
		&FollowingRepoImpl{d},
		&ArticleRepoImpl{d},
		&ArticleTagsRepository{d},
		&ArticleFavoritesRepoImpl{d},
		&CommentRepoImpl{d},
	}
}
