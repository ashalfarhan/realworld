CREATE TABLE IF NOT EXISTS comments (
    id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    comment_id      SERIAL,
	body          TEXT NOT NULL,
	created_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	updated_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	author_id         UUID NOT NULL,
	article_id         UUID NOT NULL,
    CONSTRAINT fk_comments_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_comments_article
        FOREIGN KEY (article_id)
        REFERENCES articles(id)
        ON DELETE CASCADE
);
