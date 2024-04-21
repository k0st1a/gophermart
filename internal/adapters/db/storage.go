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

	err := d.pool.QueryRow(ctx, "INSERT INTO users (login,password) VALUES($1,$2) RETURNING id", login, password).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("failed to create user:%w", err)
	}

	return id, nil
}

func (d *db) GetUser(ctx context.Context, login, password string) (int64, error) {
	log.Printf("GetUser, login:%s, password:%s", login, password)
	var id int64

	err := d.pool.QueryRow(ctx, "SELECT id FROM users WHERE login = $1 AND password = $2 LIMIT 1", login, password).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return id, ports.ErrUserNotFound
	}

	return id, nil
}

func (d *db) Close() {
	d.pool.Close()
}
