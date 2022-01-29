CREATE TABLE users (
    id           UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    user_id     SERIAL NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
	username    VARCHAR(255) UNIQUE NOT NULL,
	bio         TEXT, 
	image       VARCHAR(255),
);

CREATE TABLE followers (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    following_id UUID NOT NULL,
    follower_id UUID NOT NULL
);

CREATE TABLE articles (
    id           UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    article_id      SERIAL,
    slug           VARCHAR(255) UNIQUE NOT NULL,
	title           VARCHAR(255) UNIQUE NOT NULL,
	description     TEXT,
	body          TEXT NOT NULL,
	tag_list        VARCHAR(255)[] NOT NULL,
	created_at      TIMESTAMP DEFAULT NOW(),
	updated_at      TIMESTAMP DEFAULT NOW(),
	favorites_count   INT DEFAULT 0,
	author_id         UUID NOT NULL,
    CONSTRAINT fk_articles_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_articles_tags
        FOREIGN KEY (tag_list)
        REFERENCES tags(name)
        ON DELETE CASCADE
);

CREATE TABLE comments (
    id          UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    body        TEXT,
    author_id    UUID NOT NULL,
    CONSTRAINT fk_comments_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE tags (
    name VARCHAR(255) NOT NULL PRIMARY KEY,
);
