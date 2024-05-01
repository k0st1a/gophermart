package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

type server struct {
	server *http.Server
}

func New(ctx context.Context, address string, handler http.Handler) *server {
	s := &http.Server{
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Addr:        address,
		Handler:     handler,
	}

	return &server{
		server: s,
	}
}

func (s *server) Run() error {
	log.Printf("Run api")

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve error:%w", err)
	}

	return nil
}
