package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (d *db) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction error:%w", err)
	}

	return tx, nil
}

func (d *db) Rollback(ctx context.Context, tx pgx.Tx) error {
	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("rollback error:%w", err)
	}

	return nil
}

func (d *db) Commit(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit error:%w", err)
	}

	return nil
}
