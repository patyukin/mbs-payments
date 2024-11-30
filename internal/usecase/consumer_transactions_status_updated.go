package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) UpdateTransactionsStatus(ctx context.Context, record *kgo.Record) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		var transactions []model.TransactionSendStatus
		if err := json.Unmarshal(record.Value, &transactions); err != nil {
			return fmt.Errorf("failed to unmarshal message for topic '%s': %w", record.Topic, err)
		}

		err := repo.UpdateTransactionsSendStatus(ctx, transactions)
		if err != nil {
			return fmt.Errorf("failed repo.UpdateTransactionsStatus: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed processing in PaymentConsumeHandler: %w", err)
	}

	return nil
}
