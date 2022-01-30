CREATE TABLE IF NOT EXISTS users (
    id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id     SERIAL NOT NULL,
    email       VARCHAR(255) UNIQUE NOT NULL,
	username    VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
	bio         TEXT, 
	image       VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS followings (
    following_id UUID NOT NULL,
    follower_id UUID NOT NULL,
    PRIMARY KEY(following_id, follower_id),
    CONSTRAINT fk_followings_following
        FOREIGN KEY (following_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_followings_follower
        FOREIGN KEY (follower_id)
        REFERENCES users(id) ON DELETE CASCADE
);