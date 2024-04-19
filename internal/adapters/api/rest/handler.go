package rest

import (
	"github.com/k0st1a/gophermart/internal/ports"
)

type handler struct {
	storage ports.UserStorage
}

func newHandler(s ports.UserStorage) *handler {
	return &handler{
		storage: s,
	}
}
