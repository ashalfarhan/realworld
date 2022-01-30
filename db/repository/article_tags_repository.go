package repository

import "database/sql"

type ArticleTagsRepository struct {
	db *sql.DB
}

func (r *ArticleTagsRepository) InsertOne(articleID string, tagName string) error {
	_, err := r.db.Exec(`
	INSERT INTO 
		article_tags
		(article_id, tag_name)
	VALUES
		($1, $2)
	`, articleID, tagName)

	return err
}

func (r *ArticleTagsRepository) GetArticleTagsById(articleID string) ([]string, error) {
	row, err := r.db.Query(`
	SELECT 
		tag_name
	WHERE
		article_id = $1
	`, articleID)

	if err != nil {
		return nil, err
	}

	defer row.Close()
	var tags []string
	for row.Next() {
		var tag string
		if err := row.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
