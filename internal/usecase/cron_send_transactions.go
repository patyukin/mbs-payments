package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) CronSendTransactions(ctx context.Context) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		transactions, err := repo.SelectNotSentTransactions(ctx)
		if err != nil {
			return fmt.Errorf("failed repo.SelectNotSentTransactions: %w", err)
		}

		if len(transactions) == 0 {
			log.Info().Msg("no not sent transactions")
			return nil
		}

		err = repo.UpdateNotSentTransactions(ctx, transactions)
		if err != nil {
			return fmt.Errorf("failed repo.UpdateNotSentTransactions: %w", err)
		}

		transactionsBytes, err := json.Marshal(transactions)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.kfk.PublishTransactionReport(ctx, transactionsBytes)
		if err != nil {
			return fmt.Errorf("failed u.kfk.PublishPaymentReport: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return nil
}
