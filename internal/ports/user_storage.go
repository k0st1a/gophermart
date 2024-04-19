package ports

import "errors"

type UserStorage interface {
	RegisterUser(login, password string) error
}

var (
	ErrLoginAlreadyBusy = errors.New("login is already busy")
)
