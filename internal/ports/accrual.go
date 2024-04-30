package ports

import (
	"context"
	"errors"
)

var (
	ErrOrderNotRegistered = errors.New("order not registered")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrBlocked            = errors.New("blocked")
)

type Getter interface {
	Get(ctx context.Context, order string) (*Accrual, error)
}

type Accrual struct {
	Order   string
	Status  string
	Accrual float64
}
