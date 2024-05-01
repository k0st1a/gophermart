package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/k0st1a/gophermart/internal/ports"
)

type Managment interface {
	Create(ctx context.Context, login, password string) (int64, error)
	GetIDAndPassword(ctx context.Context, login string) (int64, string, error)
	GetBalanceAndWithdrawn(ctx context.Context, userID int64) (float64, float64, error)
}

type user struct {
	storage ports.UserStorage
}

var (
	ErrLoginAlreadyBusy = errors.New("user login is already busy")
	ErrNotFound         = errors.New("user not found")
)

func New(storage ports.UserStorage) Managment {
	return &user{
		storage: storage,
	}
}

func (u *user) Create(ctx context.Context, login, password string) (int64, error) {
	id, err := u.storage.CreateUser(ctx, login, password)
	if err != nil {
		if errors.Is(err, ports.ErrLoginAlreadyBusy) {
			return 0, ErrLoginAlreadyBusy
		}

		return 0, fmt.Errorf("storage error of create user:%w", err)
	}

	return id, nil
}

func (u *user) GetIDAndPassword(ctx context.Context, login string) (int64, string, error) {
	id, password, err := u.storage.GetUserIDAndPassword(ctx, login)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return id, password, ErrNotFound
		}

		return id, password, fmt.Errorf("storage error of get user id and password:%w", err)
	}

	return id, password, nil
}

func (u *user) GetBalanceAndWithdrawn(ctx context.Context, userID int64) (float64, float64, error) {
	current, withdrawn, err := u.storage.GetBalanceAndWithdrawn(ctx, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("storage error of get balance:%w", err)
	}

	return current, withdrawn, nil
}
