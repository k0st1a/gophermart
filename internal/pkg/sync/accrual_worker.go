package accrual

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type Accrual struct {
	Order   int64
	Status  string
	Accrual float64
}

type worker struct {
	client  ports.Getter
	accrual chan Accrual
	order   <-chan int64
}

func NewAccrualWorker(client ports.Getter, order <-chan int64) (*worker, <-chan Accrual) {
	accrual := make(chan Accrual)
	return &worker{
		client:  client,
		order:   order,
		accrual: accrual,
	}, accrual
}

func (w *worker) Run(ctx context.Context) error {
	log.Printf("Run accrual worker")

	var orderID int64

	for {
		select {
		case <-ctx.Done():
			log.Printf("Accrual worker closed with cause:%s", ctx.Err())
			return nil
		case orderID = <-w.order:
			log.Printf("Got order %d from channel", orderID)

			order := strconv.FormatInt(orderID, 10)

			apiAccrual, err := w.client.Get(ctx, order)
			if errors.Is(err, ports.ErrOrderNotRegistered) {
				log.Printf("For orderId:%v, order not registered", orderID)
				continue
			}
			if err != nil {
				return fmt.Errorf("client error of get accrual for order:%w", err)
			}

			log.Printf("For orderID:%v, apiAccrual:%+v", orderID, apiAccrual)

			orderID, err := strconv.ParseInt(apiAccrual.Order, 10, 64)
			if err != nil {
				return fmt.Errorf("strconv error of parse accrual order:%w", err)
			}

			accrual := Accrual{
				Order:   orderID,
				Status:  apiAccrual.Status,
				Accrual: apiAccrual.Accrual,
			}

			log.Printf("For orderID:%v, accrual:%v", orderID, accrual)

			select {
			case w.accrual <- accrual:
			case <-ctx.Done():
				log.Printf("Accrual worker closed with cause:%s", ctx.Err())
				return nil
			}
		}
	}
}
