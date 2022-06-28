package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/model"
	"github.com/jmoiron/sqlx"
)

type CommentRepoImpl struct {
	db *sqlx.DB
}

type CommentRepository interface {
	InsertOne(context.Context, *model.Comment) error
	FindByArticleID(context.Context, string) ([]*model.Comment, error)
	DeleteByID(context.Context, string) error
	FindOneByID(context.Context, string) (*model.Comment, error)
}

func (r *CommentRepoImpl) InsertOne(ctx context.Context, c *model.Comment) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO article_comments (body, author_username, article_id)
	VALUES (:body, :author_username, :article_id)
	RETURNING id, created_at, updated_at`
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err = stmt.GetContext(ctx, c, *c); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *CommentRepoImpl) FindByArticleID(ctx context.Context, articleID string) ([]*model.Comment, error) {
	var comments []*model.Comment
	query := `
	SELECT 
		ac.id, ac.body, ac.created_at, ac.updated_at,
		us.username as "author.username", us.bio AS "author.bio", us.image AS "author.image"
	FROM article_comments AS ac 
	LEFT JOIN users AS us 
		ON us.username = ac.author_username
	WHERE ac.article_id = $1
	ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &comments, query, articleID); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepoImpl) DeleteByID(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "DELETE FROM article_comments as ac WHERE ac.id = $1"
	if _, err = tx.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *CommentRepoImpl) FindOneByID(ctx context.Context, id string) (*model.Comment, error) {
	comm := &model.Comment{}
	query := `
	SELECT id, body, author_username, created_at, updated_at
	FROM article_comments as ac WHERE ac.id = $1`
	if err := r.db.GetContext(ctx, comm, query, id); err != nil {
		return nil, err
	}
	return comm, nil
}
