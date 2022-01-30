CREATE TABLE IF NOT EXISTS article_tags (
    article_id UUID NOT NULL,
    tag_name VARCHAR(255) NOT NULL,
    PRIMARY KEY(article_id, tag_name),
    CONSTRAINT fk_article
        FOREIGN KEY(article_id)
        REFERENCES articles(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_tag
        FOREIGN KEY(tag_name)
        REFERENCES tags(name)
        ON DELETE CASCADE
);