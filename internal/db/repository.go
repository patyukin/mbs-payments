package db

import (
	"context"
	"database/sql"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, q string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...interface{}) *sql.Row
}

type Repository struct {
	db QueryExecutor
}
