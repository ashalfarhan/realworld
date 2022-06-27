CREATE TABLE IF NOT EXISTS followings (
    following_username    VARCHAR(255) NOT NULL,
    follower_username     VARCHAR(255) NOT NULL,
    PRIMARY KEY(following_username, follower_username),
    CONSTRAINT fk_followings_following
        FOREIGN KEY (following_username)
        REFERENCES users(username) ON DELETE CASCADE,
    CONSTRAINT fk_followings_follower
        FOREIGN KEY (follower_username)
        REFERENCES users(username) ON DELETE CASCADE
);
