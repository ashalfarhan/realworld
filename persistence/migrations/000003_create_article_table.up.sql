CREATE TABLE IF NOT EXISTS articles (
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    article_id      SERIAL,
    slug            VARCHAR(255) UNIQUE NOT NULL,
	title           VARCHAR(255) NOT NULL,
	description     TEXT NOT NULL,
	body            TEXT NOT NULL,
	created_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	updated_at      TIMESTAMP DEFAULT NOW() NOT NULL,
	author_username VARCHAR(255) NOT NULL,
    CONSTRAINT fk_articles_author
        FOREIGN KEY (author_username)
        REFERENCES users(username)
        ON DELETE CASCADE
);
