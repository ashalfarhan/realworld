package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func (r *CommentRepository) InsertOne(ctx context.Context, c *model.Comment) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO
		comments
		(body, author_id, article_id)
	VALUES
		(:body, :author_id, :article_id)
	RETURNING
		id, created_at, updated_at`
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

type FindCommentsByArticleIDArgs struct {
	ArticleID string `db:"article_id"`
	Limit     int    `validate:"min=1,max=25" db:"limit"`
	Offset    int    `validate:"min=0" db:"offset"`
}

func (r *CommentRepository) FindByArticleID(ctx context.Context, args *FindCommentsByArticleIDArgs) ([]*model.Comment, error) {
	var comments []*model.Comment
	query := `
	SELECT
		id, body, author_id, created_at, updated_at
	FROM 
		comments
	WHERE 
		article_id = :article_id
	ORDER BY created_at DESC
	LIMIT :limit
	OFFSET :offset`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.SelectContext(ctx, &comments, args); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) DeleteByID(ctx context.Context, commentID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `
	DELETE FROM
		comments
	WHERE
		comments.id = $1`
	if _, err = tx.ExecContext(ctx, query, commentID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CommentRepository) FindOneByID(ctx context.Context, commentID string) (*model.Comment, error) {
	var comm *model.Comment
	query := `
	SELECT
		id, body, author_id, created_at, updated_at
	FROM
		comments
	WHERE
		comments.author_id = $1`
	if err := r.db.GetContext(ctx, comm, query, commentID); err != nil {
		return nil, err
	}

	return comm, nil
}
