-- name: GetToys :many
SELECT * FROM toys
WHERE available >= sqlc.arg(min_available)
LIMIT ? OFFSET ?;

-- name: GetToysByIdentifier :many
SELECT * FROM toys
WHERE available >= sqlc.arg(min_available) AND identifier in sqlc.slice(ids);

-- name: GetToyByID :one
SELECT * FROM toys
WHERE id = sqlc.arg(id)
LIMIT 1;

-- name: AddToyToCatalog :one
INSERT INTO toys(identifier,available,json_data)
VALUES (?,0,?)
RETURNING id;

-- name: TakeToyFromInventory :exec
UPDATE toys
SET available = available - sqlc.arg(amount)
WHERE id = sqlc.narg(id) OR identifier = sqlc.narg(identifier)
LIMIT 1;
