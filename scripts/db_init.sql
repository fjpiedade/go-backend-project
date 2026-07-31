CREATE DATABASE social

CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    update_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE posts (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255),
    user_id VARCHAR(255),
    tags TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    update_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_post_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

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
        'Primeiro Post',
        'Olá mundo!',
        1,
        'go,backend'
    );