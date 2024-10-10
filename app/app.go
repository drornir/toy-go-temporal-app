package app

import (
	"database/sql"

	"github.com/drornir/toy-go-temporal-app/pkg/sql/sqlc"
	"github.com/drornir/toy-go-temporal-app/pkg/toys"
)

func New(db *sql.DB) *App {
	return &App{
		db: db,
		shop: &toys.Shop{
			DB:   db,
			Repo: sqlc.New(db),
		},
	}
}

// /
type App struct {
	db   *sql.DB
	shop *toys.Shop
}
