CREATE TABLE IF NOT EXISTS users (
    id           UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id     SERIAL,
    email       VARCHAR(255) UNIQUE NOT NULL,
	username    VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
	bio         TEXT,
	image       VARCHAR(255),
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
