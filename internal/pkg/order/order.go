package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/k0st1a/gophermart/internal/ports"
)

type OrderManagment interface {
	CreateOrder(ctx context.Context, userID, orderID int64) error
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
