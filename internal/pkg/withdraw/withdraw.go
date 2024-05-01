package withdraw

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type Managment interface {
	Create(ctx context.Context, userID, orderID int64, sum float64) error
	List(ctx context.Context, userID int64) ([]Withdraw, error)
}

type Withdraw struct {
	Order       int64
	Sum         float64
	ProcessedAt time.Time
}

var ErrNotEnoughFunds = errors.New("not enough funds in balance")

type withdraw struct {
	storage ports.WithdrawStorage
}

func New(storage ports.WithdrawStorage) Managment {
	return &withdraw{
		storage: storage,
	}
}

func (w *withdraw) Create(ctx context.Context, userID, orderID int64, sum float64) error {
	log.Printf("Create withdraw, userID:%v, orderID:%v, sum:%v", userID, orderID, sum)

	tx, err := w.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("storage error of begin transition:%w", err)
	}
	defer func() {
		_ = w.storage.Commit(ctx, tx)
	}()

	balance, withdraw, err := w.storage.GetBalanceAndWithdrawnWithBlock(ctx, tx, userID)
	if err != nil {
		_ = w.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get balance:%w", err)
	}
	log.Printf("For userID:%v, balance:%v, withdraw:%v", userID, balance, withdraw)

	if balance < sum {
		log.Printf("For userID:%v, not enough balance", userID)
		_ = w.storage.Rollback(ctx, tx)
		return ErrNotEnoughFunds
	}

	err = w.storage.UpdateBalanceAndWithdrawn(ctx, tx, userID, balance-sum, withdraw+sum)
	if err != nil {
		_ = w.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of update balance:%w", err)
	}

	err = w.storage.CreateWithdraw(ctx, tx, userID, orderID, sum)
	if err != nil {
		_ = w.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of create withdraw:%w", err)
	}

	return nil
}

func (w *withdraw) List(ctx context.Context, userID int64) ([]Withdraw, error) {
	log.Printf("Get list of withdrawals, userID:%v", userID)
	var withdrawals []Withdraw
	dbWithdrawals, err := w.storage.GetWithdrawals(ctx, userID)
	if err != nil {
		return withdrawals, fmt.Errorf("storage error of get withdrawals:%w", err)
	}
	for _, dbWithdraw := range dbWithdrawals {
		processedAt, err := time.Parse(time.RFC3339, dbWithdraw.ProcessedAt.Format(time.RFC3339))
		if err != nil {
			return withdrawals, fmt.Errorf("error of parse ProcessedAt to RFC3339:%w", err)
		}

		withdrawals = append(withdrawals, Withdraw{
			Order:       dbWithdraw.Order,
			Sum:         dbWithdraw.Sum,
			ProcessedAt: processedAt,
		})
	}

	return withdrawals, nil
}
