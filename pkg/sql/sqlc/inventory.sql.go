// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: inventory.sql

package sqlc

import (
	"context"
	"database/sql"
	"strings"
)

const addToyToCatalog = `-- name: AddToyToCatalog :one
INSERT INTO toys(identifier,available,json_data)
VALUES (?,0,?)
RETURNING id
`

type AddToyToCatalogParams struct {
	Identifier string         `db:"identifier" json:"identifier"`
	JsonData   sql.NullString `db:"json_data" json:"json_data"`
}

func (q *Queries) AddToyToCatalog(ctx context.Context, arg AddToyToCatalogParams) (int64, error) {
	row := q.queryRow(ctx, q.addToyToCatalogStmt, addToyToCatalog, arg.Identifier, arg.JsonData)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getToyByID = `-- name: GetToyByID :one
SELECT id, identifier, available, json_data FROM toys
WHERE id = ?1
LIMIT 1
`

func (q *Queries) GetToyByID(ctx context.Context, id int64) (Toy, error) {
	row := q.queryRow(ctx, q.getToyByIDStmt, getToyByID, id)
	var i Toy
	err := row.Scan(
		&i.ID,
		&i.Identifier,
		&i.Available,
		&i.JsonData,
	)
	return i, err
}

const getToys = `-- name: GetToys :many
SELECT id, identifier, available, json_data FROM toys
WHERE available >= ?3
LIMIT ? OFFSET ?
`

type GetToysParams struct {
	MinAvailable int64 `db:"min_available" json:"min_available"`
	Limit        int64 `db:"limit" json:"limit"`
	Offset       int64 `db:"offset" json:"offset"`
}

func (q *Queries) GetToys(ctx context.Context, arg GetToysParams) ([]Toy, error) {
	rows, err := q.query(ctx, q.getToysStmt, getToys, arg.MinAvailable, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Toy
	for rows.Next() {
		var i Toy
		if err := rows.Scan(
			&i.ID,
			&i.Identifier,
			&i.Available,
			&i.JsonData,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getToysByIdentifier = `-- name: GetToysByIdentifier :many
SELECT id, identifier, available, json_data FROM toys
WHERE available >= ?1 AND identifier in /*SLICE:idents*/?
`

type GetToysByIdentifierParams struct {
	MinAvailable int64    `db:"min_available" json:"min_available"`
	Idents       []string `db:"idents" json:"idents"`
}

func (q *Queries) GetToysByIdentifier(ctx context.Context, arg GetToysByIdentifierParams) ([]Toy, error) {
	query := getToysByIdentifier
	var queryParams []interface{}
	queryParams = append(queryParams, arg.MinAvailable)
	if len(arg.Idents) > 0 {
		for _, v := range arg.Idents {
			queryParams = append(queryParams, v)
		}
		query = strings.Replace(query, "/*SLICE:idents*/?", strings.Repeat(",?", len(arg.Idents))[1:], 1)
	} else {
		query = strings.Replace(query, "/*SLICE:idents*/?", "NULL", 1)
	}
	rows, err := q.query(ctx, nil, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Toy
	for rows.Next() {
		var i Toy
		if err := rows.Scan(
			&i.ID,
			&i.Identifier,
			&i.Available,
			&i.JsonData,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const takeToyFromInventory = `-- name: TakeToyFromInventory :exec
UPDATE toys
SET available = available - ?1
WHERE id = ?2 OR identifier = ?3
LIMIT 1
`

type TakeToyFromInventoryParams struct {
	Amount     int64          `db:"amount" json:"amount"`
	ID         sql.NullInt64  `db:"id" json:"id"`
	Identifier sql.NullString `db:"identifier" json:"identifier"`
}

func (q *Queries) TakeToyFromInventory(ctx context.Context, arg TakeToyFromInventoryParams) error {
	_, err := q.exec(ctx, q.takeToyFromInventoryStmt, takeToyFromInventory, arg.Amount, arg.ID, arg.Identifier)
	return err
}
