package ports

import (
	"context"
	"database/sql"
	"errors"
	"time"
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
