package application

import (
	"context"
	"fmt"

	"github.com/k0st1a/gophermart/internal/adapters/api/accrual"
	"github.com/k0st1a/gophermart/internal/adapters/api/rest"
	"github.com/k0st1a/gophermart/internal/adapters/db"
	"github.com/k0st1a/gophermart/internal/pkg/auth"
	"github.com/k0st1a/gophermart/internal/pkg/order"
	job "github.com/k0st1a/gophermart/internal/pkg/sync"
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
	defer db.Close()

	auth := auth.New(cfg.SecretKey)
	user := user.New(db)
	order := order.New(db)
	withdraw := withdraw.New(db)

	h := rest.NewHandler(auth, user, order, withdraw)
	r := rest.BuildRoute(h, auth)

	a := accrual.New(cfg.AccrualSystemAddress)

	op, orderCh := job.NewOrderPoller(1, db)
	aw, accrualCh := job.NewAccrualWorker(a, orderCh)
	ou := job.NewOrderUpdater(db, accrualCh)

	go func() {
		err := op.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error of run order poller")
		}
	}()
	go func() {
		err := aw.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error of run accrual worker")
		}
	}()
	go func() {
		err := ou.Run(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error of run order updater")
		}
	}()

	server := rest.New(ctx, cfg.RunAddress, r)
	err = server.Run()
	if err != nil {
		return fmt.Errorf("failed to run server:%w", err)
	}

	return nil
}
