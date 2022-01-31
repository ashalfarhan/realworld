package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ashalfarhan/realworld/conduit"
	"github.com/ashalfarhan/realworld/db/model"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) InsertOne(u *model.User) error {
	return r.db.
		QueryRow(`
	INSERT INTO
		users
		(email, username, password, bio, image)
	VALUES
		($1, $2, $3, $4, $5)
	RETURNING
		users.id, users.bio, users.image`,
			u.Email,
			u.Username,
			u.Password,
			u.Bio,
			u.Image,
		).Scan(&u.ID, &u.Bio, &u.Image)
}

func (r *UserRepository) FindOneById(id string, u *model.User) error {
	return r.db.
		QueryRow(`
	SELECT
		id, email, username, bio, image, created_at, updated_at
	FROM
		users
	WHERE
		users.id = $1`, id).
		Scan(&u.ID, &u.Email, &u.Username, &u.Bio, &u.Image, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindOne(cand *model.User) error {
	return r.db.
		QueryRow(`
	SELECT
		id, email, username, password, bio, image 
	FROM
		users 
	WHERE
		users.email = $1 
	OR
		users.username = $2`,
			cand.Email,
			cand.Username,
		).
		Scan(&cand.ID, &cand.Email, &cand.Username, &cand.Password, &cand.Bio, &cand.Image)
}

func (r *UserRepository) UpdateOne(u *conduit.UpdateUserArgs) error {
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

	if  u.Bio.Set {
		argIdx++
		updateArgs = append(updateArgs, fmt.Sprintf("bio = $%d", argIdx))
		valArgs = append(valArgs, u.Bio)
	}

	if  u.Image.Set {
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

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	if _, err := stmt.Exec(valArgs...); err != nil {
		return err
	}

	return nil
}
