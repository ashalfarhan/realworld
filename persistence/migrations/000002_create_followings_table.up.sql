CREATE TABLE IF NOT EXISTS followings (
    following_id    UUID NOT NULL,
    follower_id     UUID NOT NULL,
    PRIMARY KEY(following_id, follower_id),
    CONSTRAINT fk_followings_following
        FOREIGN KEY (following_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_followings_follower
        FOREIGN KEY (follower_id)
        REFERENCES users(id) ON DELETE CASCADE
);
