-- name: CreateUserRole :one
INSERT INTO user_roles (role_id, user_id)
VALUES ($1, $2) RETURNING *;

-- name: ListUserRoles :many
SELECT * FROM user_roles;

-- name: GetUsersByRole :many
SELECT * FROM user_roles WHERE role_id = $1;

-- name: UpdateUserRole :one
UPDATE user_roles
SET role_id = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserRole :exec
DELETE FROM user_roles
WHERE id = $1;


