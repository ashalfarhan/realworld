CREATE TABLE IF NOT EXISTS article_favorites (
    username    VARCHAR(255) NOT NULL,
    article_id  UUID NOT NULL,
    PRIMARY KEY(username, article_id),
    CONSTRAINT fk_article_favorites_user
        FOREIGN KEY (username)
        REFERENCES users(username) ON DELETE CASCADE,
    CONSTRAINT fk_article_favorites_article
        FOREIGN KEY (article_id)
        REFERENCES articles(id) ON DELETE CASCADE
);