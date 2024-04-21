package rest

import (
// "github.com/k0st1a/gophermart/internal/ports"
)

type handler struct {
	storage any
}

func newHandler(s any) *handler {
	return &handler{
		storage: s,
	}
}
