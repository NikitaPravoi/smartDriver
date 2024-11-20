-- name: CreateUser :one
INSERT INTO users (login, password, name, surname, patronymic, organization_id)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetOrganizationUsers :many
SELECT * FROM users WHERE organization_id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, surname = $3, password = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserWithHisRoles :one
SELECT u.id, u.login, r.name AS role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1;

-- name: GetUserWithHisRolesByOrganization :one
SELECT u.id, u.login, r.name AS role_name
FROM users u
         JOIN user_roles ur ON u.id = ur.user_id
         JOIN roles r ON ur.role_id = r.id
WHERE u.organization_id = $1;


