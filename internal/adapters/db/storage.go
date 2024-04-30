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

func (d *db) GetOrderUserID(ctx context.Context, tx pgx.Tx, orderID int64) (int64, error) {
	log.Printf("GetOrderUserID, orderID:%v", orderID)
	var userID int64

	err := d.pool.QueryRow(ctx, "SELECT user_id FROM orders WHERE id = $1", orderID).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ports.ErrOrderNotFound
	}

	return userID, nil
}

func (d *db) CreateOrder(ctx context.Context, tx pgx.Tx, userID, orderID int64) error {
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

func (d *db) GetOrders(ctx context.Context, tx pgx.Tx, userID int64) ([]ports.Order, error) {
	var orders []ports.Order

	rows, err := tx.Query(ctx,
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

func (d *db) Close() {
	d.pool.Close()
}
