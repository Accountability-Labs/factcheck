-- name: CreateNote :one
INSERT INTO notes (created_by, created_at, updated_at, note, url)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

-- name: UpdateNote :one
UPDATE notes
SET note = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetRecentNNotesForUrl :many
SELECT * FROM notes WHERE url = $1 ORDER BY created_at DESC LIMIT $2;