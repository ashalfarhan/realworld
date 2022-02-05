package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ashalfarhan/realworld/db/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// See https://go.dev/doc/database/execute-transactions
func (r *UserRepository) InsertOne(ctx context.Context, u *model.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Defer a rollback incase returning error

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
		return err
	}

	defer stmt.Close()
	if err = stmt.GetContext(ctx, u, *u); err != nil {
		return err
	}

	// Commit and then return the error if any
	return tx.Commit()
}

func (r *UserRepository) FindOneByID(ctx context.Context, id string, u *model.User) error {
	query := `
	SELECT
		id, email, username, bio, image, created_at, updated_at
	FROM
		users
	WHERE
		users.id = $1`
	return r.db.GetContext(ctx, u, query, id)

}

func (r *UserRepository) FindOne(ctx context.Context, cand *model.User) error {
	query := `
	SELECT
		id, email, username, password, bio, image 
	FROM
		users 
	WHERE
		users.email = $1 
	OR
		users.username = $2`
	return r.db.GetContext(ctx, cand, query, cand.Email, cand.Username)
}

// Use Pointer to update
// To determine if the field needs to be updated
type UpdateUserValues struct {
	ID       string
	Email    *string
	Username *string
	Password *string
	Image    model.NullString
	Bio      model.NullString
}

func (r *UserRepository) UpdateOne(ctx context.Context, u *UpdateUserValues) error {
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
	valArgs = append(valArgs, u.ID)
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
