package application

import (
	"context"
	"fmt"

	"github.com/k0st1a/gophermart/internal/adapters/api/rest"
	"github.com/k0st1a/gophermart/internal/adapters/db"
	"github.com/rs/zerolog/log"
)

func Run() error {
	log.Debug().Msg("Running application")
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

	api := rest.NewAPI(ctx, cfg.RunAddress, db)
	err = api.Run()
	if err != nil {
		return fmt.Errorf("failed to run api:%w", err)
	}

	return nil
}
