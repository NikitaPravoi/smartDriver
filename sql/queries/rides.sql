-- name: CreateRide :one
INSERT INTO rides (
    branch_id,
    created_at
) VALUES (
             $1,
             CURRENT_TIMESTAMP
         )
RETURNING *;

-- name: GetRide :one
SELECT * FROM rides
WHERE id = $1;

-- name: CompleteRide :exec
UPDATE rides
SET ended_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: AttachOrderToRide :exec
INSERT INTO rides_to_orders (
    ride_id,
    order_id
) VALUES (
             $1, $2
         );

-- name: DetachAllOrdersFromRide :exec
DELETE FROM rides_to_orders
WHERE ride_id = $1;

-- name: GetOrdersByRideID :many
SELECT o.*
FROM orders o
         JOIN rides_to_orders rto ON rto.order_id = o.id
WHERE rto.ride_id = $1
ORDER BY o.created_at;

-- name: GetActiveRides :many
SELECT r.*,
       COUNT(rto.order_id) as order_count
FROM rides r
         LEFT JOIN rides_to_orders rto ON r.id = rto.ride_id
WHERE r.ended_at IS NULL
GROUP BY r.id
ORDER BY r.created_at DESC;