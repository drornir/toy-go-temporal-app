package toys

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/drornir/toy-go-temporal-app/pkg/sql/sqlc"
)

type (
	ShopOrderForm struct {
		CustomerID uint
		Items      []ShopOrderFormItem
	}
	ShopOrderFormItem struct {
		ToyIdentifier string
		Amount        uint
	}
	OrderReciept struct {
		OrderNumber   uint64
		OriginalOrder ShopOrderForm
	}
	ItemAvailability struct {
		ToyIdentifier string
		Available     uint
	}
)

type Shop struct {
	DB   *sql.DB
	Repo *sqlc.Queries
}

func (self *Shop) ReserveOrderFromInventory(ctx context.Context, order ShopOrderForm) error {
	idents := make([]string, len(order.Items))
	for idx, item := range order.Items {
		idents[idx] = item.ToyIdentifier
	}
	dbToys, err := self.Repo.GetToysByIdentifier(ctx, sqlc.GetToysByIdentifierParams{Idents: idents})
	if err != nil {
		return fmt.Errorf("getting toys by ids from db: %w", err)
	}
	idents = nil

	itemsMap := make(map[string]string, len(dbToys))
	toyToAmount := make(map[string]int64, len(dbToys))
	for _, dbToy := range dbToys {
		for _, orderItem := range order.Items {
			if dbToy.Identifier != orderItem.ToyIdentifier {
				continue
			}
			if orderItem.Amount > uint(dbToy.Available) {
				return fmt.Errorf("can't fullfil order for toy %q: requested %d but only %d are available", orderItem.ToyIdentifier, orderItem.Amount, dbToy.Available)
			}
			itemAsJson, err := json.Marshal(orderItem)
			if err != nil {
				return fmt.Errorf("serializing orderItem to JSON: for Go value %#v: error: %w", orderItem, err)
			}
			itemsMap[dbToy.Identifier] = string(itemAsJson)
			toyToAmount[dbToy.Identifier] = int64(orderItem.Amount)
		}
	}

	tx, err := self.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begining db tansaction: %w", err)
	}
	repoWithTx := self.Repo.WithTx(tx)

	for toyIdent, amount := range toyToAmount {
		err := repoWithTx.TakeToyFromInventory(ctx, sqlc.TakeToyFromInventoryParams{
			Amount:     amount,
			Identifier: sql.NullString{Valid: true, String: toyIdent},
		})
		if err != nil {
			err := errors.Join(err, tx.Rollback())
			return fmt.Errorf("taking %d %q toys from inventory: %w", amount, toyIdent, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting inventory changes to db: %w", err)
	}

	return nil
}

func (self *Shop) CreateReceipt(ctx context.Context, order ShopOrderForm) (OrderReciept, error) {
	tx, err := self.DB.BeginTx(ctx, nil)
	if err != nil {
		return OrderReciept{}, fmt.Errorf("begining db tansaction: %w", err)
	}
	repoWithTx := self.Repo.WithTx(tx)
	jsonDataForOrder, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}
	sqlOrder, err := repoWithTx.CreateOrder(ctx, sqlc.CreateOrderParams{
		CustomerID: int64(order.CustomerID),
		JsonData:   string(jsonDataForOrder),
	})
	if err != nil {
		tx.Rollback()
		return OrderReciept{}, fmt.Errorf("creating order in db: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return OrderReciept{}, fmt.Errorf("commiting order and inventory changes to db: %w", err)
	}

	return OrderReciept{
		OrderNumber:   uint64(sqlOrder.ID),
		OriginalOrder: order,
	}, nil
}
