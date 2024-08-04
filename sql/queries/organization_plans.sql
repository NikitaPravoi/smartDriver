-- name: CreateOrganizationPlan :one
INSERT INTO organization_plans (organization_id, plan_id, start_date, end_date)
    VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListOrganizationPlans :many
SELECT * FROM organization_plans;

-- name: GetOrganizationPlan :one
SELECT * FROM organization_plans WHERE id = $1;

-- name: UpdateOrganizationPlan :exec
UPDATE organization_plans
SET organization_id = $2, plan_id = $3, start_date = $4, end_date = $5
WHERE id = $1;

-- name: DeleteOrganizationPlan :exec
DELETE FROM organization_plans
WHERE id = $1;
