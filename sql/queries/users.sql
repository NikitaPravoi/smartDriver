-- name: CheckUserLoginExists :one
SELECT EXISTS(
    SELECT 1 FROM users WHERE login = $1
) AS exists;

-- name: CreateUser :one
INSERT INTO users (
    login,
    password,
    name,
    surname,
    patronymic,
    organization_id,
    created_at,
    updated_at
) VALUES (
             $1, $2, $3, $4, $5, $6,
             CURRENT_TIMESTAMP,
             CURRENT_TIMESTAMP
         )
RETURNING *;

-- name: GetUser :one
SELECT u.*, array_agg(r.name) as roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1
GROUP BY u.id;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE($2, name),
    surname = COALESCE($3, surname),
    patronymic = COALESCE($4, patronymic),
    password = COALESCE($5, password),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListUsers :many
SELECT u.*, array_agg(r.name) as roles
FROM users u
         LEFT JOIN user_roles ur ON u.id = ur.user_id
         LEFT JOIN roles r ON ur.role_id = r.id
WHERE
    ($1::bigint IS NULL OR u.organization_id = $1) AND
    (
        $2 = '' OR
        u.login ILIKE $2 OR
        u.name ILIKE $2 OR
        u.surname ILIKE $2
        )
GROUP BY u.id
ORDER BY u.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUsers :one
SELECT COUNT(DISTINCT u.id)
FROM users u
WHERE
    ($1::bigint IS NULL OR u.organization_id = $1) AND
    (
        $2 = '' OR
        u.login ILIKE $2 OR
        u.name ILIKE $2 OR
        u.surname ILIKE $2
        );

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;

-- name: DeleteUserRoles :exec
DELETE FROM user_roles
WHERE user_id = $1;

-- name: AssignUserRole :exec
INSERT INTO user_roles (
    user_id,
    role_id
) VALUES (
             $1, $2
         );

-- name: GetUserRoles :many
SELECT r.*
FROM roles r
         JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = CURRENT_TIMESTAMP
WHERE id = $1;