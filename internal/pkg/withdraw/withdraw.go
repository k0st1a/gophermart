package withdraw

import (
	"context"
	"errors"
	"fmt"

	"github.com/k0st1a/gophermart/internal/ports"
)

type Managment interface {
	Create(ctx context.Context, userID, orderID int64, sum float64) error
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
	tx, err := w.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("storage error of begin transition:%w", err)
	}
	defer func() {
		_ = w.storage.Commit(ctx, tx)
	}()

	balance, withdraw, err := w.storage.GetBalanceWithBlock(ctx, tx, userID)
	if err != nil {
		_ = w.storage.Rollback(ctx, tx)
		return fmt.Errorf("storage error of get balance:%w", err)
	}

	if balance < sum {
		_ = w.storage.Rollback(ctx, tx)
		return ErrNotEnoughFunds
	}

	err = w.storage.UpdateBalance(ctx, tx, userID, balance-sum, withdraw+sum)
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
