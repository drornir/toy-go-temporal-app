-- name: CreateCustomer :one
INSERT INTO customers(name) VALUES (sqlc.narg(name)) RETURNING *;

-- name: CreateOrder :one
INSERT INTO orders(customer_id, json_data) VALUES (?, ?) RETURNING *;
