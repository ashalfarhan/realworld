package repository

import (
	"context"

	"github.com/ashalfarhan/realworld/api/dto"
	"github.com/ashalfarhan/realworld/db/model"
	"github.com/jmoiron/sqlx"
)

type UserRepoImpl struct {
	db *sqlx.DB
}

type UserRepository interface {
	InsertOne(context.Context, *dto.RegisterUserFields) (*model.User, error)
	FindOneByID(context.Context, string) (*model.User, error)
	FindOne(context.Context, *FindOneUserFilter) (*model.User, error)
	UpdateOne(context.Context, *dto.UpdateUserFields, *model.User) error
}

// See https://go.dev/doc/database/execute-transactions
func (r *UserRepoImpl) InsertOne(ctx context.Context, d *dto.RegisterUserFields) (*model.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Defer a rollback incase returning error
	u := &model.User{
		Email:    d.Email,
		Username: d.Username,
		Password: d.Password,
	}

	query := `
	INSERT INTO
		users
		(email, username, password)
	VALUES
		(:email, :username, :password)
	RETURNING
		users.id, users.bio, users.image`
	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err = stmt.GetContext(ctx, u, *d); err != nil {
		return nil, err
	}

	// Commit and then return the error if any
	return u, tx.Commit()
}

func (r *UserRepoImpl) FindOneByID(ctx context.Context, id string) (*model.User, error) {
	u := &model.User{
		ID: id,
	}

	query := `
	SELECT
		id, email, username, bio, image, created_at, updated_at
	FROM
		users
	WHERE
		users.id = $1`

	if err := r.db.GetContext(ctx, u, query, id); err != nil {
		return nil, err
	}
	return u, nil

}

type FindOneUserFilter struct {
	Email    string
	Username string
}

func (r *UserRepoImpl) FindOne(ctx context.Context, d *FindOneUserFilter) (*model.User, error) {
	u := new(model.User)

	query := `
	SELECT
		id, email, username, password, bio, image 
	FROM
		users 
	WHERE
		users.email = $1 
	OR
		users.username = $2`

	if err := r.db.GetContext(ctx, u, query, d.Email, d.Username); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepoImpl) UpdateOne(ctx context.Context, d *dto.UpdateUserFields, u *model.User) error {
	if v := d.Email; v != nil {
		u.Email = *v
	}
	if v := d.Username; v != nil {
		u.Username = *v
	}
	if v := d.Password; v != nil {
		u.Password = *v
	}

	if v := d.Bio; v.Set {
		u.Bio = v
	}
	if v := d.Image; v.Set {
		u.Image = v
	}

	query := `
	UPDATE users
	SET
		email = :email,
		username = :username,
		password = :password,
		bio = :bio,
		image = :image,
		updated_at = NOW()
	WHERE
		users.id = :id
	`
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, u); err != nil {
		return err
	}

	return tx.Commit()
}
