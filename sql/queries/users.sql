-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, api_key, email, user_name, password_hash, salt)
VALUES (NOW(), NOW(), encode(sha256(random()::text::bytea), 'hex'), $1, $2, $3, $4)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: GetSaltByEmail :one
SELECT salt FROM users WHERE email = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;