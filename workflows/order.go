package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/drornir/toy-go-temporal-app/pkg/toys"
)

func RegisterOrderWorkflow(w worker.Worker) {
	var shop *toys.Shop
	w.RegisterWorkflow(Order)
	w.RegisterActivity(shop.ReserveOrderFromInventory)
	w.RegisterActivity(shop.CreateReceipt)
}

type OrderWorkflowParams struct {
	Order toys.ShopOrderForm
}

func Order(ctx workflow.Context, params OrderWorkflowParams) (toys.OrderReciept, error) {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var shop *toys.Shop
	reserverActivity := workflow.ExecuteActivity(ctx, shop.ReserveOrderFromInventory, params.Order)
	if err := reserverActivity.Get(ctx, nil); err != nil {
		return toys.OrderReciept{}, fmt.Errorf("reserving order from inventory: %w", err)
	}

	var receipt toys.OrderReciept
	receiptActivity := workflow.ExecuteActivity(ctx, shop.CreateReceipt, params.Order)
	if err := receiptActivity.Get(ctx, &receipt); err != nil {
		return toys.OrderReciept{}, fmt.Errorf("creating order receipt: %w", err)
	}

	return receipt, nil
}
