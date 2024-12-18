package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/patyukin/bs-payments/internal/model"
	pkgModel "github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *Repository) InsertTransaction(ctx context.Context, in model.Payment, accountID, typeTransaction string) (string, error) {
	query := `
INSERT INTO transactions (payment_id, account_id, type, amount, currency, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
`
	row := r.db.QueryRowContext(ctx, query, in.ID, accountID, typeTransaction, in.Amount, in.Currency, "COMPLETED", time.Now().UTC())
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

func (r *Repository) SelectNotSentTransactions(ctx context.Context) ([]pkgModel.Transaction, error) {
	query := `
SELECT 
    t.id,
    t.payment_id, 
    t.account_id,
    a.user_id,
    t.type,
    t.amount,
    t.currency,
    t.description,
    p.description AS payment_description,
    t.status,
    t.send_status,
    t.created_at 
FROM transactions t
INNER JOIN accounts a ON t.account_id = a.id
INNER JOIN payments p ON t.payment_id = p.id
WHERE t.send_status = 'PENDING'
ORDER BY t.id;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed during row iteration in rows.Err(): %w", err)
	}

	var transactions []pkgModel.Transaction
	for rows.Next() {
		var transaction pkgModel.Transaction
		if err = rows.Scan(
			&transaction.ID,
			&transaction.PaymentID,
			&transaction.AccountID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Description,
			&transaction.PaymentDescription,
			&transaction.Status,
			&transaction.SendStatus,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *Repository) UpdateNotSentTransactions(ctx context.Context, transactions []pkgModel.Transaction) error {
	query := `UPDATE transactions SET send_status = $1 WHERE id = $2`
	for _, transaction := range transactions {
		_, err := r.db.ExecContext(ctx, query, "PROCESSING", transaction.ID)
		if err != nil {
			return fmt.Errorf("failed r.db.ExecContext: %w", err)
		}
	}

	return nil
}

func (r *Repository) UpdateTransactionsSendStatus(ctx context.Context, transactions []pkgModel.TransactionSendStatus) error {
	query := `UPDATE transactions SET send_status = $1 WHERE id = $2`
	for _, transaction := range transactions {
		_, err := r.db.ExecContext(ctx, query, transaction.SendStatus, transaction.ID)
		if err != nil {
			return fmt.Errorf("failed r.db.ExecContext: %w", err)
		}
	}

	return nil
}

func (r *Repository) SelectTransactionsByPaymentID(ctx context.Context, paymentID string) ([]model.Transaction, error) {
	query := `
SELECT 
		id,
		payment_id, 
		account_id,
		type,
		amount,
		currency,
		description,
		status,
		send_status,
		created_at 
FROM transactions WHERE payment_id = $1
ORDER BY id
`
	rows, err := r.db.QueryContext(ctx, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed r.db.QueryContext: %w", err)
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Error().Msgf("failed rows.Close: %v", err)
		}
	}(rows)

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed during row iteration in rows.Err(): %w", err)
	}

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		if err = rows.Scan(
			&transaction.ID,
			&transaction.PaymentID,
			&transaction.AccountID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Description,
			&transaction.Status,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed rows.Scan: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
