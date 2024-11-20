-- name: CreateBranch :one
INSERT INTO branches (name, location, organization_id)
    VALUES ($1, $2, $3) RETURNING *;

-- name: ListBranches :many
SELECT * FROM branches;

-- name: GetBranch :one
SELECT * FROM branches WHERE id = $1;

-- name: UpdateBranch :one
UPDATE branches
SET name = $2, location = $3, organization_id = $4
WHERE id = $1
RETURNING *;

-- name: DeleteBranch :exec
DELETE FROM branches
WHERE id = $1;
