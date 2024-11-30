package db

import (
	"context"
	"fmt"
	"github.com/patyukin/bs-payments/internal/model"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/rs/zerolog/log"
	"time"
)

func (r *Repository) InsertPayment(ctx context.Context, in *paymentpb.CreatePaymentRequest) (string, error) {
	log.Debug().Msgf("in: %v", in)

	currentTime := time.Now().UTC()
	query := `
INSERT INTO payments (sender_account_id, receiver_account_id, amount, currency, description, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.SenderAccountId, in.ReceiverAccountId, in.Amount, in.Currency, in.Description, "DRAFT", currentTime)
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

func (r *Repository) GetPaymentByID(ctx context.Context, paymentID string) (model.Payment, error) {
	query := `
SELECT
	id,
	sender_account_id, 
	receiver_account_id,
	amount,
	currency,
	description,
	status,
	created_at,
	updated_at
FROM payments WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, paymentID)

	if row.Err() != nil {
		return model.Payment{}, fmt.Errorf("failed r.db.QueryRowContext: %w", row.Err())
	}

	var payment model.Payment
	err := row.Scan(
		&payment.ID,
		&payment.SenderAccountID,
		&payment.ReceiverAccountID,
		&payment.Amount,
		&payment.Currency,
		&payment.Description,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return model.Payment{}, fmt.Errorf("failed row.Scan: %w", err)
	}

	return payment, nil
}

func (r *Repository) UpdatePaymentStatusByID(ctx context.Context, paymentID, status string) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, paymentID)
	if err != nil {
		return fmt.Errorf("failed r.db.ExecContext: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed result.RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
