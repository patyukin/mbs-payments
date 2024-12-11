package db

import (
	"context"
	"fmt"
	"time"

	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (r *Repository) InsertAccount(ctx context.Context, in *paymentpb.CreateAccountRequest) (string, error) {
	query := `INSERT INTO accounts (user_id, currency, balance, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.UserId, in.Currency, in.Balance, time.Now().UTC())
	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext, row.Err(): %w", row.Err())
	}

	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

func (r *Repository) SelectBankTypeAccountID(ctx context.Context) (string, error) {
	query := `SELECT id FROM accounts WHERE type = 'BANK'`
	row := r.db.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var id string
	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return id, nil
}

func (r *Repository) ExistsPositiveAccountBalance(ctx context.Context, accountID string) (int64, error) {
	query := `SELECT balance FROM accounts WHERE id = $1 AND balance > 0`
	row := r.db.QueryRowContext(ctx, query, accountID)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var balance int64
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed row.Scan: %w", err)
	}

	return balance, nil
}

func (r *Repository) UpdateAccountBalance(ctx context.Context, accountID string, balance int64) error {
	query := `UPDATE accounts SET balance = balance - $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, accountID, balance)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) IncreaseAccountBalance(ctx context.Context, accountID string, balance int64) error {
	query := `UPDATE accounts SET balance = balance + $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, accountID, balance)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) DecreaseAccountBalance(ctx context.Context, accountID string, balance int64) error {
	query := `UPDATE accounts SET balance = balance - $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, accountID, balance)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	return nil
}

func (r *Repository) SelectUserIDByAccountID(ctx context.Context, accountID string) (string, error) {
	query := `SELECT user_id FROM accounts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, accountID)
	if row.Err() != nil {
		return "", fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("failed row.Scan: %w", err)
	}

	return userID, nil
}
