CREATE DATABASE social;

-- extension to postgres
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- inverted operation
DROP TABLE IF EXISTS users;

-- ##########################################

CREATE TABLE IF NOT EXISTS posts (
    -- id BIGSERIAL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255),
    user_id BIGINT,
    tags TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- other way
CREATE TABLE IF NOT EXISTS posts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ALTER TABLE if not created the FK when table created

-- add on up
ALTER TABLE posts
ADD CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users (id);

-- add on down
ALTER TABLE posts DROP CONSTRAINT fk_post_user;

--up
ALTER TABLE posts
ALTER COLUMN user_id TYPE BIGINT
USING user_id::BIGINT;

ALTER TABLE posts
ADD CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

-- down
ALTER TABLE posts DROP CONSTRAINT fk_post_user;

ALTER TABLE posts ALTER COLUMN user_id TYPE VARCHAR(255);

-- inserts

INSERT INTO
    users (
        first_name,
        last_name,
        username,
        email,
        password
    )
VALUES (
        'Ensei',
        'Tankado',
        'ensei.tankado',
        'ensei.tankado@email.com',
        '123456'
    );

INSERT INTO
    posts (title, content, user_id, tags)
VALUES (
        'First Post',
        'Hello World!',
        1,
        'go,backend'
    );