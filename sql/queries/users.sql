-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, api_key, email)
VALUES (NOW(), NOW(), encode(sha256(random()::text::bytea), 'hex'), $1)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE api_key = $1;