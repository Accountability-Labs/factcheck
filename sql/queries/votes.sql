-- name: VoteOnNote :one
INSERT INTO votes (voted_by, voted_on, created_at, vote)
VALUES ($1, $2, NOW(), $3)
RETURNING *;

-- name: UpdateVoteOnNote :one
UPDATE votes
SET vote = $2, voted_on = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteVote :exec
DELETE FROM votes WHERE id = $1;