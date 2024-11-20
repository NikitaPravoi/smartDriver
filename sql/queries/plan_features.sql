-- name: CreatePlanFeature :one
INSERT INTO plan_features (plan_id, feature_name)
    VALUES ($1, $2) RETURNING *;

-- name: ListPlanFeatures :many
SELECT * FROM plan_features;

-- name: GetPlanFeature :one
SELECT * FROM plan_features WHERE id = $1;

-- name: UpdatePlanFeature :one
UPDATE plan_features
SET plan_id = $2, feature_name = $3
WHERE id = $1
RETURNING *;

-- name: DeletePlanFeature :exec
DELETE FROM plan_features
WHERE id = $1;
