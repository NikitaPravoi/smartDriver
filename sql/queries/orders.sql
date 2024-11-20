-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (
    customer_name,
    phone,
    city,
    street,
    apartment,
    floor,
    doorphone,
    building,
    entrance,
    comment,
    cost,
    status,
    location,
    created_at,
    external_id
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, point($13, $14), $15, $16
         )
RETURNING *;

-- name: UpdateOrder :exec
UPDATE orders
SET
    customer_name = $2,
    phone = $3,
    city = $4,
    street = $5,
    apartment = $6,
    floor = $7,
    doorphone = $8,
    building = $9,
    entrance = $10,
    comment = $11,
    cost = $12,
    location = point($13, $14)
WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1
ORDER BY created_at DESC;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) FROM orders
WHERE status = $1;

-- name: GetOrderByExternalID :one
SELECT * FROM orders
WHERE external_id = $1;