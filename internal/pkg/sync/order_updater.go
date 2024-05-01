package accrual

import (
	"context"
	"fmt"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type updater struct {
	storage ports.OrderUpdaterStorage
	accrual <-chan Accrual
}

func NewOrderUpdater(storage ports.OrderUpdaterStorage, accrual <-chan Accrual) *updater {
	return &updater{
		storage: storage,
		accrual: accrual,
	}
}

func (u *updater) Run(ctx context.Context) error {
	log.Printf("Run order updater")

	var accrual Accrual

	for {
		select {
		case <-ctx.Done():
			log.Printf("Order updater closed with cause:%s", ctx.Err())
			return nil
		case accrual = <-u.accrual:
			log.Printf("Get accrual %+v from channel", accrual)

			err := u.updateOrder(ctx, accrual.Order, accrual.Status, accrual.Accrual)
			if err != nil {
				return fmt.Errorf("storage error of update order:%w", err)
			}
		}
	}
}

func (u *updater) updateOrder(ctx context.Context, orderID int64, status string, accrual float64) error {
	tx, err := u.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("storage error of begin transition:%w", err)
	}

	userID, err := u.storage.GetUserIDByOrderWithBlock(ctx, tx, orderID)
	if err != nil {
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get user id by order:%w", err)
	}

	balance, err := u.storage.GetBalanceWithBlock(ctx, tx, userID)
	if err != nil {
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get balance with block:%w", err)
	}

	err = u.storage.UpdateOrder(ctx, tx, orderID, status, accrual)
	if err != nil {
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of update order:%w", err)
	}

	if accrual != 0 {
		err = u.storage.UpdateBalance(ctx, tx, userID, balance+accrual)
		if err != nil {
			_ = u.storage.Rollback(ctx, tx)
			return fmt.Errorf("storage error of update balance:%w", err)
		}
	}

	err = u.storage.Commit(ctx, tx)
	if err != nil {
		return fmt.Errorf("storage error of commit transaction:%w", err)
	}

	return nil
}
