-- name: CreateUserRole :one
INSERT INTO user_roles (role_id, user_id)
VALUES ($1, $2) RETURNING *;

-- name: ListUserRoles :many
SELECT * FROM user_roles;

-- name: GetUserRoles :many
SELECT * FROM user_roles WHERE user_id = $1;

-- name: GetUsersByRole :many
SELECT * FROM user_roles WHERE role_id = $1;

-- name: UpdateUserRole :exec
UPDATE user_roles
SET role_id = $2
WHERE id = $1;

-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE id = $1;

