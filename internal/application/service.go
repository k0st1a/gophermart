package application

import (
	"context"
	"fmt"

	"github.com/k0st1a/gophermart/internal/adapters/api/rest"
	"github.com/k0st1a/gophermart/internal/adapters/db"
	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/k0st1a/gophermart/internal/pkg/order"
	"github.com/k0st1a/gophermart/internal/pkg/user"
	"github.com/k0st1a/gophermart/internal/pkg/withdraw"
	"github.com/rs/zerolog/log"
)

func Run() error {
	log.Printf("Running application")
	ctx := context.Background()

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	db, err := db.NewDB(ctx, cfg.DatabaseURI)
	if err != nil {
		return fmt.Errorf("failed to create db:%w", err)
	}

	auth := auth.New(cfg.SecretKey)
	user := user.New(db)
	order := order.New(db)
	withdraw := withdraw.New(db)

	h := rest.NewHandler(auth, user, order, withdraw)
	r := rest.BuildRoute(h, auth)

	server := rest.New(ctx, cfg.RunAddress, r)
	err = server.Run()
	if err != nil {
		return fmt.Errorf("failed to run server:%w", err)
	}

	return nil
}
