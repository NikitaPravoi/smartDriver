-- name: CreateRideToOrder :one
INSERT INTO rides_to_orders (ride_id, order_id)
    VALUES ($1, $2) RETURNING *;

-- name: ListRidesToOrders :many
SELECT * FROM rides_to_orders;

-- name: GetRideToOrder :one
SELECT * FROM rides_to_orders WHERE id = $1;

-- name: UpdateRideToOrder :one
UPDATE rides_to_orders
SET ride_id = $2, order_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteRideToOrder :exec
DELETE FROM rides_to_orders
WHERE id = $1;

