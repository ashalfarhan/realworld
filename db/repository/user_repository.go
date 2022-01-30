package repository

import (
	"database/sql"

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
		id, email, username, bio, image
	FROM 
		users
	WHERE 
		users.id = $1`, id).
		Scan(&u.ID, &u.Email, &u.Username, &u.Bio, &u.Image)
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

func (d *UserRepository) UpdateOne(u *model.User) error {
	_, err := d.db.
		Exec(`
	UPDATE 
		users
	SET 
		email = $2,
		username = $3,
		password = $4,
		bio = $5,
		image = $6,
		updated_at = NOW()
	WHERE users.id = $1`,
			u.ID,
			u.Email,
			u.Username,
			u.Password,
			u.Bio,
			u.Image,
		)
	return err
}
