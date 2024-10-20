package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) *sql.Row
}

type Client struct {
	db *sql.DB
}

func (c *Client) GetRepo() *Repository {
	return &Repository{
		db: c.db,
	}
}

type Handler func(ctx context.Context, repo *Repository) error

func New(db *sql.DB) *Client {
	return &Client{db: db}
}

func (c *Client) ReadCommitted(ctx context.Context, f Handler) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Error().Msgf("failed to rollback transaction: %v", errRollback)
			}
		}
	}()

	repo := &Repository{db: tx}

	if err = f(ctx, repo); err != nil {
		return fmt.Errorf("failed to execute handler: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.db.Close()
}
