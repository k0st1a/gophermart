package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/k0st1a/gophermart/internal/ports"
)

type Managment interface {
	Create(ctx context.Context, userID, orderID int64) error
	List(ctx context.Context, userID int64) ([]Order, error)
}

type Order struct {
	UploadedAt time.Time
	Status     string
	Number     int64
	Accrual    float64
}

var (
	ErrAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrAlreadyUploadedByThisUser    = errors.New("order already uploaded by this user")
)

type order struct {
	storage ports.OrderStorage
}

func New(storage ports.OrderStorage) Managment {
	return &order{
		storage: storage,
	}
}

func (o *order) Create(ctx context.Context, userID, orderID int64) error {
	dbUserID, err := o.storage.GetUserIDByOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, ports.ErrOrderNotFound) {
			err = o.storage.CreateOrder(ctx, userID, orderID)
			if err != nil {
				return fmt.Errorf("failed to create order:%w", err)
			}

			return nil
		}

		return fmt.Errorf("error of get order:%w", err)
	}

	if dbUserID != userID {
		return ErrAlreadyUploadedByAnotherUser
	}

	return ErrAlreadyUploadedByThisUser
}

func (o *order) List(ctx context.Context, userID int64) ([]Order, error) {
	orders := []Order{}

	dbOrders, err := o.storage.GetOrders(ctx, userID)
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
