-- +goose Up

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
        encode(sha256(random()::text::bytea), 'hex')
    ),
    user_name TEXT NOT NULL,
    email TEXT NOT NULL
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    note TEXT NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE votes (
    voted_by INTEGER NOT NULL REFERENCES users(id),
    voted_on INTEGER NOT NULL REFERENCES notes(id),
    PRIMARY KEY (voted_by, voted_on),
    created_at TIMESTAMP NOT NULL,
    vote INTEGER NOT NULL
);

-- +goose Down

DROP TABLE users;
DROP TABLE notes;
DROP TABLE votes;