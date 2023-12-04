-- name: GetStats :one
SELECT
( SELECT COUNT(*) FROM users) AS num_users,
( SELECT COUNT(*) FROM notes) AS num_notes,
( SELECT COUNT(*) FROM votes) AS num_votes;
