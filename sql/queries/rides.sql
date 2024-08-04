-- name: CreateRide :one
INSERT INTO rides (branch_id)
    VALUES ($1) RETURNING *;

-- name: ListRides :many
SELECT * FROM rides;

-- name: GetRide :one
SELECT * FROM rides WHERE id = $1;

-- name: DeleteRide :exec
DELETE FROM rides
WHERE id = $1;
