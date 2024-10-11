package app

import (
	"context"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/drornir/toy-go-temporal-app/workflows"
)

func (a *App) RunAsTemporalWorker(ctx context.Context, client client.Client, taskQueue string, options worker.Options) error {
	worker := worker.New(client, taskQueue, options)

	workflows.RegisterOrderWorkflow(worker)

	interruptCh := make(chan interface{}, 1)
	go func() {
		<-ctx.Done()
		interruptCh <- struct{}{}
		close(interruptCh)
	}()

	return worker.Run(interruptCh)
}
