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
SELECT users.user_name, votes.vote, notes.* FROM notes
JOIN users ON notes.created_by = users.id
LEFT JOIN votes ON (votes.voted_on = notes.id AND votes.voted_by = $1)
WHERE url = $2 ORDER BY users.created_at DESC LIMIT $3;