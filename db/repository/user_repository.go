package repository

import (
	"context"
	"fmt"
	"strings"

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
	UpdateOne(context.Context, *dto.UpdateUserFields, string) error
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
		email, username, bio, image, created_at, updated_at
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

func (r *UserRepoImpl) UpdateOne(ctx context.Context, u *dto.UpdateUserFields, uid string) error {
	var updateArgs []string
	var valArgs []interface{}
	argIdx := 0

	if v := u.Email; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("email = $%d", argIdx))
		valArgs = append(valArgs, *u.Email)
	}

	if v := u.Username; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("username = $%d", argIdx))
		valArgs = append(valArgs, *u.Username)
	}

	if u.Bio.Set {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("bio = $%d", argIdx))
		valArgs = append(valArgs, u.Bio)
	}

	if u.Image.Set {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("image = $%d", argIdx))
		valArgs = append(valArgs, u.Image)
	}

	if v := u.Password; v != nil {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("password = $%d", argIdx))
		valArgs = append(valArgs, *u.Password)
	}

	updateArgs = append(updateArgs, "updated_at = NOW()")

	argIdx++
	valArgs = append(valArgs, uid)
	query := fmt.Sprintf("UPDATE users SET %s WHERE users.id = $%d", strings.Join(updateArgs, ", "), argIdx)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, valArgs...); err != nil {
		return err
	}

	return tx.Commit()
}
