-- name: CreatePlan :one
INSERT INTO plans (name, cost, employee_limit)
    VALUES ($1, $2, $3) RETURNING *;

-- name: ListPlans :many
SELECT * FROM plans;

-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1;

-- name: UpdatePlan :one
UPDATE plans
SET name = $2, cost = $3, employee_limit = $4
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans
WHERE id = $1;
