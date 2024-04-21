package application

import (
	"context"

	"github.com/k0st1a/gophermart/internal/adapter/api/rest"
	"github.com/k0st1a/gophermart/internal/adapter/db"
	"github.com/rs/zerolog/log"
)

func Run() error {
	ctx := context.Background()

	cfg, err := collectConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	db := db.NewDB(ctx, cfg.DatabaseURI)

	api := rest.NewAPI(ctx, cfg.RunAddress, db)
	err = api.Run()
	if err != nil {
		return err
	}

	return nil
}
