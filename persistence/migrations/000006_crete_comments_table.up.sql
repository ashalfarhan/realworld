CREATE TABLE IF NOT EXISTS article_comments (
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    comment_id      SERIAL,
	body            TEXT NOT NULL,
	created_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	updated_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	author_username VARCHAR(255) NOT NULL,
	article_id      UUID NOT NULL,
    CONSTRAINT fk_article_comments_author
        FOREIGN KEY (author_username)
        REFERENCES users(username)
        ON DELETE CASCADE,
    CONSTRAINT fk_article_comments_article
        FOREIGN KEY (article_id)
        REFERENCES articles(id)
        ON DELETE CASCADE
);
