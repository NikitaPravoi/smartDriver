-- name: CreateSession :one
INSERT INTO sessions (user_id, session_token, refresh_token, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: ListSessions :many
SELECT * FROM sessions;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1;

-- name: GetSessionByToken :one
SELECT * FROM sessions WHERE session_token = $1;

-- name: UpdateSession :exec
UPDATE sessions
SET user_id = $2, session_token = $3, refresh_token = $4, created_at = $5, expires_at = $6
WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: UpdateSessionExpiry :exec
UPDATE sessions
SET expires_at = CURRENT_TIMESTAMP
WHERE refresh_token = $1;