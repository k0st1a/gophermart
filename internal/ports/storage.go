package ports

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error)
	GetBalanceAndWithdrawn(ctx context.Context, userID int64) (float64, float64, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login is already busy")
	ErrUserNotFound     = errors.New("user not found")
)

type OrderStorage interface {
	GetUserIDByOrder(ctx context.Context, orderID int64) (int64, error)
	CreateOrder(ctx context.Context, userID, orderID int64) error
	GetOrders(ctx context.Context, userID int64) ([]Order, error)
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	Status     string
	UploadedAt time.Time
	Accrual    sql.NullFloat64
	Number     int64
}

type WithdrawStorage interface {
	CreateWithdraw(ctx context.Context, tx pgx.Tx, userID, orderID int64, sum float64) error
	GetBalanceAndWithdrawnWithBlock(ctx context.Context, tx pgx.Tx, userID int64) (float64, float64, error)
	UpdateBalanceAndWithdrawn(ctx context.Context, tx pgx.Tx, userID int64, balance, withdrawn float64) error
	GetWithdrawals(ctx context.Context, userID int64) ([]Withdraw, error)

	BeginTx(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

type Withdraw struct {
	Order       int64
	Sum         float64
	ProcessedAt time.Time
}

type OrderPollerStorage interface {
	GetNotProcessedOrders(ctx context.Context) ([]int64, error)
}

type OrderUpdaterStorage interface {
	GetUserIDByOrderWithBlock(ctx context.Context, tx pgx.Tx, orderID int64) (int64, error)
	GetBalanceWithBlock(ctx context.Context, tx pgx.Tx, userID int64) (float64, error)
	UpdateOrder(ctx context.Context, tx pgx.Tx, orderID int64, status string, accrual float64) error
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, balance float64) error

	BeginTx(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
