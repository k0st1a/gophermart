package accrual

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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
			log.Printf("Got accrual %+v from channel", accrual)

			err := u.updateOrder(ctx, accrual.Order, accrual.Status, accrual.Accrual)
			if err != nil {
				log.Error().Err(err).Msg("storage error of update order")
				continue
			}
		}
	}
}

func (u *updater) updateOrder(ctx context.Context, orderID int64, status string, accrual float64) error {
	log.Printf("updateOrder, orderID:%v, status:%v, accrual:%v", orderID, status, accrual)

	tx, err := u.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("storage error of begin transition:%w", err)
	}
	defer func() {
		err := u.storage.Commit(ctx, tx)
		if errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("transaction closed for update order, orderID:%v", orderID)
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("storage error of commit transaction")
		}
	}()

	userID, err := u.storage.GetUserIDByOrderWithBlock(ctx, tx, orderID)
	if err != nil {
		log.Error().Err(err).Msg("storage error of get user id by order")
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get user id by order:%w", err)
	}
	log.Printf("For orderID:%v, userID:%v", orderID, userID)

	balance, err := u.storage.GetBalanceWithBlock(ctx, tx, userID)
	if err != nil {
		log.Error().Err(err).Msg("storage error of get balance with block")
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get balance with block:%w", err)
	}
	log.Printf("For userID:%v, balance:%v", orderID, balance)

	err = u.storage.UpdateOrder(ctx, tx, orderID, status, accrual)
	if err != nil {
		log.Error().Err(err).Msg("storage error of update order")
		_ = u.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of update order:%w", err)
	}

	if accrual != 0 {
		log.Printf("accrual not 0 => update balance, userID:%v, new balance:%v", userID, balance+accrual)
		err = u.storage.UpdateBalance(ctx, tx, userID, balance+accrual)
		if err != nil {
			log.Error().Err(err).Msg("storage error of update balance")
			_ = u.storage.Rollback(ctx, tx)
			return fmt.Errorf("storage error of update balance:%w", err)
		}
	}

	return nil
}
