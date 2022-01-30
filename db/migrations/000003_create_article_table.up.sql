CREATE TABLE IF NOT EXISTS articles (
    id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    article_id      SERIAL,
    slug           VARCHAR(255) UNIQUE NOT NULL,
	title           VARCHAR(255) NOT NULL,
	description     TEXT,
	body          TEXT NOT NULL,
	created_at      TIMESTAMP DEFAULT NOW(),
	updated_at      TIMESTAMP DEFAULT NOW(),
	author_id         UUID NOT NULL,
    CONSTRAINT fk_articles_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);