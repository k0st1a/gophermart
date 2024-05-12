package sync

import (
	"context"
	"time"

	"github.com/k0st1a/gophermart/internal/adapters/api/accrual"
	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type tick struct {
	orderStorage   ports.NotProcessedOrderStorage
	accrualAddress string
	pollInterval   int
}

func NewTicker(address string, interval int, storage ports.NotProcessedOrderStorage) *tick {
	return &tick{
		accrualAddress: address,
		pollInterval:   interval,
		orderStorage:   storage,
	}
}

func (t *tick) Run(ctx context.Context) error {
	log.Printf("Run ticker")
	ticker := time.NewTicker(time.Duration(t.pollInterval) * time.Second)

	tick := 0

	block := accrual.NewBlock()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Ticker closed with cause:%s", ctx.Err())
			ticker.Stop()
			return nil
		case <-ticker.C:
			tick++
			log.Printf("Got tick %d", tick)
			a := accrual.New(t.accrualAddress, block)
			j := NewJob(tick, t.orderStorage, a)
			err := j.Run(ctx)
			if err != nil {
				log.Error().Err(err).Msg("error of run job")
			}
		}
	}
}
