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
	GetBalance(ctx context.Context, userID int64) (float64, float64, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login is already busy")
	ErrUserNotFound     = errors.New("user not found")
)

type OrderStorage interface {
	GetOrderUserID(ctx context.Context, orderID int64) (int64, error)
	CreateOrder(ctx context.Context, userID, orderID int64) error
	GetOrders(ctx context.Context, userID int64) ([]Order, error)
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	Number     int64
	Status     string
	Accrual    sql.NullFloat64
	UploadedAt time.Time
}

type WithdrawStorage interface {
	CreateWithdraw(ctx context.Context, tx pgx.Tx, userID, orderID int64, sum float64) error
	GetBalanceWithBlock(ctx context.Context, tx pgx.Tx, userID int64) (float64, float64, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, balance, withdrawn float64) error

	BeginTx(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
