CREATE TABLE IF NOT EXISTS article_favorites (
    user_id     UUID NOT NULL,
    article_id  UUID NOT NULL,
    PRIMARY KEY(user_id, article_id),
    CONSTRAINT fk_article_favorites_user
        FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_article_favorites_article
        FOREIGN KEY (article_id)
        REFERENCES articles(id) ON DELETE CASCADE
);