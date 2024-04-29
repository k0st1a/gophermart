package ports

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login is already busy")
	ErrUserNotFound     = errors.New("user not found")
)

type OrderStorage interface {
	GetOrderUserID(ctx context.Context, tx pgx.Tx, orderID int64) (int64, error)
	CreateOrder(ctx context.Context, tx pgx.Tx, userID, orderID int64) error

	BeginTx(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

var (
	ErrOrderNotFound = errors.New("order not found")
)
