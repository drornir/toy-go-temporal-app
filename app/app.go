package app

import (
	"database/sql"

	"github.com/drornir/toy-go-temporal-app/pkg/sql/sqlc"
	"github.com/drornir/toy-go-temporal-app/pkg/toys"
)

func New(db *sql.DB) *App {
	return &App{
		Shop: &toys.Shop{DB: db, Repo: sqlc.New(db)},
	}
}

type App struct {
	Shop *toys.Shop

	stopChan chan struct{}
}
