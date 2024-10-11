package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/drornir/toy-go-temporal-app/app"
	"github.com/drornir/toy-go-temporal-app/pkg/sql"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error from temporal-worker: %s\n", err)
		os.Exit(42)
	}
}

func run() error {
	ctx, signalStopFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalStopFunc()

	db, err := sql.ConnectLibsqlDev()
	if err != nil {
		panic(err)
	}
	app := app.New(db)

	tClient, err := client.DialContext(ctx, client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})

	if err := app.RunAsTemporalWorker(ctx, tClient, "toy-go-temporal-app", worker.Options{}); err != nil {
		return fmt.Errorf("error running worker: %w", err)
	}

	return nil
}
