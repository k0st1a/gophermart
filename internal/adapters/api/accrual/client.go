package accrual

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

var isBlocked atomic.Bool

type client struct {
	address string
	client  *http.Client
}

func New(address string) ports.AccrualGetter {
	return &client{
		address: address,
		client:  &http.Client{},
	}
}

func (c *client) Get(ctx context.Context, order string) (*ports.Accrual, error) {
	log.Printf("Get accrual for order with number:%s", order)

	if isBlocked.Load() {
		return nil, ports.ErrBlocked
	}

	url, err := url.JoinPath(c.address, "/api/orders/", order)
	if err != nil {
		return nil, fmt.Errorf("url error of join path:%w", err)
	}

	log.Printf("For order:%s, url:%s", order, url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http error of new request with context:%w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error of client do:%w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("For order:%s, StatusOK", order)

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("io error of read body:%w", err)
		}
		log.Printf("For order:%s, data:%s", order, string(data))

		accrual, err := Deserialize(data)
		if err != nil {
			return nil, err
		}
		log.Printf("For order:%s, accrual:%v", order, accrual)

		return &ports.Accrual{
			Order:   accrual.Order,
			Status:  accrual.Status,
			Accrual: accrual.Accrual,
		}, nil

	case http.StatusNoContent:
		log.Printf("For order:%s, StatusNoContent", order)
		return nil, ports.ErrOrderNotRegistered

	case http.StatusTooManyRequests:
		log.Printf("For order:%s, too many requests", order)

		retryHeader := resp.Header.Get("Retry-After")
		log.Printf("For order:%s, retryHeader:%v", order, retryHeader)
		retryAfter, err := strconv.Atoi(retryHeader)
		log.Printf("For order:%s, retryAfter:%v", order, retryAfter)
		if err != nil {
			return nil, ports.ErrTooManyRequests
		}

		go block(time.Duration(retryAfter) * time.Second)

		return nil, ports.ErrTooManyRequests
	}

	return nil, fmt.Errorf("unknown response status code:%v", resp.StatusCode)
}

func block(wait time.Duration) {
	isBlocked.Store(true)
	log.Printf("block clients for %s seconds while too many requests", wait)
	time.Sleep(wait)
	isBlocked.Store(false)
	log.Printf("unblock clients after %s seconds", wait)
}
