package ports

import (
	"context"
	"errors"
)

type UserStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	GetUser(ctx context.Context, login, password string) (int64, error)
}

var (
	ErrLoginAlreadyBusy = errors.New("login is already busy")
	ErrUserNotFound     = errors.New("user not found")
)
