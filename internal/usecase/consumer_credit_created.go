package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) ConsumerCreditCreated(ctx context.Context, record *kgo.Record) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		var credit model.CreditCreated
		if err := json.Unmarshal(record.Value, &credit); err != nil {
			return fmt.Errorf("failed to unmarshal message for topic '%s': %w", record.Topic, err)
		}

		err := repo.DecreaseAccountBalance(ctx, credit.AccountID, credit.Amount)
		if err != nil {
			return fmt.Errorf("failed repo.DecreaseAccountBalance: %w", err)
		}

		u.rbtmq.PushCreditCreated(ctx, record.Value, record.Headers)
		// send notify "кредит успешно создан"

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed processing in PaymentConsumeHandler: %w", err)
	}
}
