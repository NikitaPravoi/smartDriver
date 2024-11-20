-- name: GetUserByLogin :one
SELECT * FROM users
WHERE login = $1;

-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    session_token,
    refresh_token,
    created_at,
    expires_at
) VALUES (
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: GetSessionByToken :one
SELECT s.*, u.*
FROM sessions s
         JOIN users u ON s.user_id = u.id
WHERE s.session_token = $1 AND s.expires_at > CURRENT_TIMESTAMP;

-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
WHERE refresh_token = $1;

-- name: UpdateSession :one
UPDATE sessions
SET
    session_token = $2,
    refresh_token = $3,
    expires_at = $4
WHERE id = $1
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: GetUserSessions :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;