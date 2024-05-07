package sync

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type job struct {
	storage ports.NotProcessedOrderStorage
	client  ports.AccrualGetter
	number  int
}

func NewJob(number int, storage ports.NotProcessedOrderStorage, accrual ports.AccrualGetter) *job {
	return &job{
		number:  number,
		storage: storage,
		client:  accrual,
	}
}

func (j *job) Run(ctx context.Context) error {
	log.Printf("Run job #%v", j.number)

	tx, err := j.storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("storage error of begin transition, error:%w", err)
	}
	defer func() {
		_ = j.storage.Rollback(ctx, tx)
	}()

	orderID, err := j.storage.GetNotProcessedOrderWithBlock(ctx, tx)
	if err != nil {
		return fmt.Errorf("storage error of get not processed order, error:%w", err)
	}

	log.Printf("Job #%v, get accrual for orderID:%v", j.number, orderID)
	orderString := strconv.FormatInt(orderID, 10)
	accrual, err := j.client.Get(ctx, orderString)
	if errors.Is(err, ports.ErrOrderNotRegistered) {
		log.Printf("Job #%v, orderID:%v not registered in accrual", j.number, orderID)

		err = j.storage.UpdateOrder(ctx, tx, orderID, "INVALID", 0)
		if err != nil {
			return fmt.Errorf("storage error of update INVALID order, error:%w", err)
		}
		err = j.storage.Commit(ctx, tx)
		if err != nil {
			return fmt.Errorf("storage error of commit transaction, error:%w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("client error of get accrual for order, error:%w", err)
	}

	log.Printf("Job #%v, for orderID:%v, accrual:%+v", j.number, orderID, accrual)

	if accrual.Order != orderString {
		log.Printf("Job #%v, other accrual order from response, order from request:%v"+
			", order from response:%v", j.number, orderString, accrual.Order)

		err = j.storage.UpdateOrder(ctx, tx, orderID, "INVALID", 0)
		if err != nil {
			return fmt.Errorf("storage error of update INVALID order, error:%w", err)
		}

		err = j.storage.Commit(ctx, tx)
		if err != nil {
			return fmt.Errorf("storage error of commit transaction, error:%w", err)
		}
		return fmt.Errorf("other accrual order from response, order from request:%v"+
			", order from response:%v", orderString, accrual.Order)
	}

	userID, err := j.storage.GetUserIDByOrderWithBlock(ctx, tx, orderID)
	if err != nil {
		return fmt.Errorf("storage error of get user id by order, error:%w", err)
	}
	log.Printf("Job #%v, for orderID:%v, userID:%v", j.number, orderID, userID)

	balance, err := j.storage.GetBalanceWithBlock(ctx, tx, userID)
	if err != nil {
		return fmt.Errorf("storage error of get balance with block, error:%w", err)
	}
	log.Printf("Job #%v, for userID:%v, balance:%v", j.number, orderID, balance)

	err = j.storage.UpdateOrder(ctx, tx, orderID, accrual.Status, accrual.Accrual)
	if err != nil {
		return fmt.Errorf("storage error of update order, error:%w", err)
	}

	if accrual.Accrual != 0 {
		log.Printf("Job #%v, accrual not 0 => update balance, userID:%v, new balance:%v",
			j.number, userID, balance+accrual.Accrual)
		err = j.storage.UpdateBalance(ctx, tx, userID, balance+accrual.Accrual)
		if err != nil {
			return fmt.Errorf("storage error of update balance, error:%w", err)
		}
	}

	err = j.storage.Commit(ctx, tx)
	if err != nil {
		return fmt.Errorf("storage error of commit transaction, error:%w", err)
	}

	return nil
}
