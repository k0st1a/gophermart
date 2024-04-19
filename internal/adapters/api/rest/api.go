package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/k0st1a/gophermart/internal/ports"
)

type api struct {
	server *http.Server
}

func NewAPI(ctx context.Context, address string, storage ports.UserStorage) *api {
	h := newHandler(storage)

	r := newRouter()
	buildRoute(r, h)

	s := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Addr:        address,
		Handler:     r,
	}

	return &api{
		server: s,
	}
}

func (a *api) Run() error {
	err := a.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
