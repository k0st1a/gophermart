package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/k0st1a/gophermart/internal/ports"
)

type OrderManagment interface {
	CreateOrder(ctx context.Context, userID, orderID int64) error
	GetOrders(ctx context.Context, userID int64) ([]Order, error)
}

type Order struct {
	Number     int64
	Status     string
	Accrual    float64
	UploadedAt time.Time
}

var (
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrOrderAlreadyUploadedByThisUser    = errors.New("order already uploaded by this user")
)

type order struct {
	storage ports.OrderStorage
}

func New(storage ports.OrderStorage) OrderManagment {
	return &order{
		storage: storage,
	}
}

func (o *order) CreateOrder(ctx context.Context, userID, orderID int64) error {
	tx, err := o.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error of begin transaction:%w", err)
	}
	defer func() {
		_ = o.storage.Commit(ctx, tx)
	}()

	dbUserID, err := o.storage.GetOrderUserID(ctx, tx, orderID)
	if err != nil {
		if errors.Is(err, ports.ErrOrderNotFound) {
			err = o.storage.CreateOrder(ctx, tx, userID, orderID)
			if err != nil {
				return fmt.Errorf("failed to create order:%w", err)
			}

			return nil
		}

		return fmt.Errorf("error of get order:%w", err)
	}

	if dbUserID != userID {
		return ErrOrderAlreadyUploadedByAnotherUser
	}

	return ErrOrderAlreadyUploadedByThisUser
}

func (o *order) GetOrders(ctx context.Context, userID int64) ([]Order, error) {
	var orders []Order

	tx, err := o.storage.BeginTx(ctx)
	if err != nil {
		return orders, fmt.Errorf("error of begin transaction:%w", err)
	}
	defer func() {
		_ = o.storage.Commit(ctx, tx)
	}()

	dbOrders, err := o.storage.GetOrders(ctx, tx, userID)
	if err != nil {
		return orders, fmt.Errorf("error of get orders from storage:%w", err)
	}

	for _, dbOrder := range dbOrders {
		uploadedAt, err := time.Parse(time.RFC3339, dbOrder.UploadedAt.Format(time.RFC3339))
		if err != nil {
			return orders, fmt.Errorf("error of parse UploadedAt to RFC3339:%w", err)
		}

		orders = append(orders, Order{
			Number:     dbOrder.Number,
			Status:     dbOrder.Status,
			Accrual:    dbOrder.Accrual.Float64,
			UploadedAt: uploadedAt,
		})
	}

	return orders, nil
}