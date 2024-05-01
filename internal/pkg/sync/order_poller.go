package accrual

import (
	"context"
	"fmt"
	"time"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type poller struct {
	storage  ports.OrderPollerStorage
	interval int
	order    chan<- int64
}

func NewOrderPoller(interval int, storage ports.OrderPollerStorage) (*poller, <-chan int64) {
	order := make(chan int64)
	return &poller{
		interval: interval,
		storage:  storage,
		order:    order,
	}, order
}

func (p *poller) Run(ctx context.Context) error {
	log.Printf("Run order poller")
	ticker := time.NewTicker(time.Duration(p.interval) * time.Second)

	tick := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("Order poller closed with cause:%s", ctx.Err())
			ticker.Stop()
			return nil
		case <-ticker.C:
			tick++
			log.Printf("Got tick %d of order polling", tick)
			//По хорошему нужна лучашая логика получения не обработанных заказов.
			//Например, вычитывать по одному с учетом последнего полученного.
			//но, т.к. времени мало, то делаем самый "топорный" вариант -
			//вычитываем все за раз.
			orders, err := p.storage.GetNotProcessedOrders(ctx)
			if err != nil {
				return fmt.Errorf("storage error of get not processed orders:%w", err)
			}

			log.Printf("For tick %d orders:%v", tick, orders)

			for _, orderID := range orders {
				p.order <- orderID
			}
		}
	}
}
