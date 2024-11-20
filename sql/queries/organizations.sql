-- name: CreateOrganization :one
INSERT INTO organizations (name, balance, iiko_api_token)
    VALUES ($1, $2, $3) RETURNING *;

-- name: ListOrganizations :many
SELECT * FROM organizations;

-- name: GetOrganization :one
SELECT * FROM organizations WHERE id = $1;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = $2, balance = $3, iiko_api_token = $4
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: GetOrganizationsApiTokens :many
SELECT iiko_api_token FROM organizations WHERE balance > 0;
