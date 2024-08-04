-- name: CreateOrder :one
INSERT INTO orders (customer_name, customer_phone, city, street, apartment, floor, entrance, comment, cost, status, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: UpdateOrder :exec
UPDATE orders
SET customer_name = $2, customer_phone = $3, city = $4, street = $5, apartment = $6, floor = $7, entrance = $8, comment = $9, cost = $10, status = $11
WHERE id = $1;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
