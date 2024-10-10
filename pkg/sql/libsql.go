package sql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/go-libsql"
)

func ConnectLibsqlDev(ctx context.Context) (*sql.DB, error) {
	dbName := "file:./bin/toys-dev.db"
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		return nil, fmt.Errorf("opening %q using libsql driver: %w", dbName, err)
	}
	return db, nil
}
