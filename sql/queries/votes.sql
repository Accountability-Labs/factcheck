-- name: VoteOnNote :one
INSERT INTO votes (voted_by, voted_on, created_at, vote)
VALUES ($1, $2, NOW(), $3)
ON CONFLICT (voted_by, voted_on) DO UPDATE SET vote = $3, created_at = NOW()
RETURNING *;

-- name: DeleteVote :exec
DELETE FROM votes WHERE (voted_by, voted_on) = ($1, $2);