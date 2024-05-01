package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k0st1a/gophermart/internal/ports"
	"github.com/rs/zerolog/log"
)

type db struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (*db, error) {
	err := runMigrations(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	return &db{
		pool: pool,
	}, nil
}

func (d *db) CreateUser(ctx context.Context, login, password string) (int64, error) {
	log.Printf("CreateUser, login:%s, password:%s", login, password)
	var id int64

	err := d.pool.QueryRow(ctx,
		"INSERT INTO users (login,password) VALUES($1,$2) "+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
		login, password).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrLoginAlreadyBusy
	}

	if err != nil {
		return id, fmt.Errorf("failed to create user:%w", err)
	}

	return id, nil
}

func (d *db) GetUserIDAndPassword(ctx context.Context, login string) (int64, string, error) {
	log.Printf("GetUserIDAndPassword, login:%s", login)
	var id int64
	var password string

	err := d.pool.QueryRow(ctx, "SELECT id, password FROM users WHERE login = $1", login).Scan(&id, &password)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", ports.ErrUserNotFound
	}

	if err != nil {
		return 0, "", fmt.Errorf("failed to get user id and password:%w", err)
	}

	return id, password, nil
}

func (d *db) GetBalanceAndWithdrawn(ctx context.Context, userID int64) (float64, float64, error) {
	log.Printf("GetBalanceAndWithdrawn, userID:%v", userID)
	var (
		balance   float64
		withdrawn float64
	)

	err := d.pool.QueryRow(ctx,
		"SELECT balance, withdrawn FROM users WHERE id = $1", userID).Scan(&balance, &withdrawn)
	if err != nil {
		return 0, 0, fmt.Errorf("query error of get balance and withdrawn:%w", err)
	}

	return balance, withdrawn, nil
}

func (d *db) GetBalanceAndWithdrawnWithBlock(ctx context.Context, tx pgx.Tx, userID int64) (float64, float64, error) {
	log.Printf("GetBalanceAndWithdrawnWithBlock, userID:%v", userID)
	var (
		balance   float64
		withdrawn float64
	)

	err := tx.QueryRow(ctx,
		"SELECT balance, withdrawn FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&balance, &withdrawn)
	if err != nil {
		return 0, 0, fmt.Errorf("query error of get balance and withdrawn with block:%w", err)
	}

	return balance, withdrawn, nil
}

func (d *db) GetBalanceWithBlock(ctx context.Context, tx pgx.Tx, userID int64) (float64, error) {
	log.Printf("GetBalanceWithBlock, userID:%v", userID)
	var balance float64

	err := tx.QueryRow(ctx,
		"SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query error of get balance with block:%w", err)
	}

	return balance, nil
}

func (d *db) UpdateBalanceAndWithdrawn(ctx context.Context, tx pgx.Tx, userID int64, balance, withdrawn float64) error {
	var id int64

	err := tx.QueryRow(ctx,
		"UPDATE ONLY users SET balance = $1, withdrawn = $2 WHERE id = $3 RETURNING id",
		balance, withdrawn, userID).Scan(&id)
	if err != nil {
		return fmt.Errorf("query error of update balance and withdrawn:%w", err)
	}

	return nil
}

func (d *db) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, balance float64) error {
	var id int64

	err := tx.QueryRow(ctx,
		"UPDATE ONLY users SET balance = $1, WHERE id = $2 RETURNING id",
		balance, userID).Scan(&id)
	if err != nil {
		return fmt.Errorf("query error of update balance:%w", err)
	}

	return nil
}

func (d *db) GetUserIDByOrder(ctx context.Context, orderID int64) (int64, error) {
	log.Printf("GetUserIDByOrder, orderID:%v", orderID)
	var userID int64

	err := d.pool.QueryRow(ctx, "SELECT user_id FROM orders WHERE id = $1", orderID).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrOrderNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("query error of get user id by order id:%w", err)
	}

	return userID, nil
}

func (d *db) GetUserIDByOrderWithBlock(ctx context.Context, tx pgx.Tx, orderID int64) (int64, error) {
	log.Printf("GetUserIDByOrderWithBlock, orderID:%v", orderID)
	var userID int64

	err := tx.QueryRow(ctx, "SELECT user_id FROM orders WHERE id = $1 FOR UPDATE", orderID).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrOrderNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("query error of get user id by order id:%w", err)
	}

	return userID, nil
}

func (d *db) CreateOrder(ctx context.Context, userID, orderID int64) error {
	log.Printf("CreateOrder, userID:%v, orderID:%v", userID, orderID)
	var id int64

	err := d.pool.QueryRow(ctx,
		"INSERT INTO orders (id,status,user_id) VALUES($1,$2,$3) RETURNING id", orderID, "NEW", userID).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create order:%w", err)
	}

	log.Printf("Created order, id:%v", id)
	return nil
}

func (d *db) UpdateOrder(ctx context.Context, tx pgx.Tx, orderID int64, status string, accrual float64) error {
	log.Printf("UpdateOrder, orderID:%v, status:%v, accrual:%v", orderID, status, accrual)
	var id int64

	err := tx.QueryRow(ctx, "UPDATE ONLY orders SET accrual = $1, status = $2 WHERE id = $3 RETURNING id",
		accrual, status, orderID).Scan(&id)
	if err != nil {
		return fmt.Errorf("query error of update order:%w", err)
	}

	log.Printf("Updated order, id:%v", id)
	return nil
}

func (d *db) GetOrders(ctx context.Context, userID int64) ([]ports.Order, error) {
	var orders []ports.Order

	rows, err := d.pool.Query(ctx,
		"SELECT id, status, accrual, uploaded_at FROM orders "+
			"WHERE user_id = $1 ORDER BY uploaded_at",
		userID)
	if err != nil {
		return orders, fmt.Errorf("query error of get orders:%w", err)
	}

	for rows.Next() {
		var o ports.Order
		err = rows.Scan(
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
		)
		if err != nil {
			return orders, fmt.Errorf("scan error of get orders:%w", err)
		}
		orders = append(orders, o)
	}

	err = rows.Err()
	if err != nil {
		return orders, fmt.Errorf("error of get orders:%w", err)
	}

	return orders, nil
}

func (d *db) GetNotProcessedOrders(ctx context.Context) ([]int64, error) {
	var orderIDList []int64

	rows, err := d.pool.Query(ctx, "SELECT id FROM orders WHERE status in ('PROCESSING', 'NEW')")
	if err != nil {
		return orderIDList, fmt.Errorf("query error of get not processed orders:%w", err)
	}

	for rows.Next() {
		var orderID int64
		err = rows.Scan(&orderID)
		if err != nil {
			return orderIDList, fmt.Errorf("scan error of get not processed order:%w", err)
		}
		orderIDList = append(orderIDList, orderID)
	}

	err = rows.Err()
	if err != nil {
		return orderIDList, fmt.Errorf("error of get not processed orders:%w", err)
	}

	return orderIDList, nil
}

func (d *db) CreateWithdraw(ctx context.Context, tx pgx.Tx, userID, orderID int64, sum float64) error {
	var id int64

	const query = `
		INSERT INTO "withdraw" (order_id, user_id, sum)
		VALUES ($1, $2, $3)
		RETURNING  "withdraw".id;
	`

	err := tx.QueryRow(ctx, "INSERT INTO withdraw (order_id, user_id, sum) VALUES ($1, $2, $3) RETURNING id",
		orderID, userID, sum).Scan(&id)
	if err != nil {
		return fmt.Errorf("query error of create withdraw:%w", err)
	}

	log.Printf("Created withdraw, id:%v", id)
	return nil
}

func (d *db) GetWithdrawals(ctx context.Context, userID int64) ([]ports.Withdraw, error) {
	log.Printf("GetWithdrawals, userID:%v", userID)
	var withdrawals []ports.Withdraw

	rows, err := d.pool.Query(ctx,
		"SELECT order_id, sum, processed_at FROM withdrawals "+
			"WHERE user_id = $1 ORDER BY processed_at",
		userID)
	if err != nil {
		return withdrawals, fmt.Errorf("query error of get withdrawals:%w", err)
	}

	for rows.Next() {
		var w ports.Withdraw
		err = rows.Scan(
			&w.Order,
			&w.Sum,
			&w.ProcessedAt,
		)
		if err != nil {
			return withdrawals, fmt.Errorf("scan error of get withdrawals:%w", err)
		}
		withdrawals = append(withdrawals, w)
	}

	err = rows.Err()
	if err != nil {
		return withdrawals, fmt.Errorf("error of get withdrawals:%w", err)
	}

	return withdrawals, nil
}

func (d *db) Close() {
	d.pool.Close()
}
