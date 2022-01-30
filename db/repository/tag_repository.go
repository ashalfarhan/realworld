package repository

import "database/sql"

type TagRepository struct {
	db *sql.DB
}

func (r *TagRepository) InsertOne(name string) error {
	_, err := r.db.Exec(`
	INSERT INTO
		tags
		(name)
	VALUES
		($1)
	`, name)

	return err
}
