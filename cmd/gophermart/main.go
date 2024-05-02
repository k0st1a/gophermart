package main

import (
	"github.com/k0st1a/gophermart/internal/application"
	"github.com/rs/zerolog/log"
)

func main() {
	err := application.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
