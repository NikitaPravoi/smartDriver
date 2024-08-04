-- name: CreateOrganization :one
INSERT INTO organizations (name, balance, iiko_id)
    VALUES ($1, $2, $3) RETURNING *;

-- name: ListOrganizations :many
SELECT * FROM organizations;

-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;

-- name: UpdateOrganization :exec
UPDATE organizations
SET name = $2, balance = $3, iiko_id = $4
WHERE id = $1;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;
