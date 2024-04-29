package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

type api struct {
	server *http.Server
}

func NewAPI(ctx context.Context, address string, handler http.Handler) *api {
	s := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Addr:        address,
		Handler:     handler,
	}

	return &api{
		server: s,
	}
}

func (a *api) Run() error {
	log.Printf("Run api")

	err := a.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
