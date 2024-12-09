package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) ConsumeCreditCreated(ctx context.Context, record *kgo.Record) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		var credit model.CreditCreated
		if err := json.Unmarshal(record.Value, &credit); err != nil {
			return fmt.Errorf("failed to unmarshal message for topic '%s': %w", record.Topic, err)
		}

		err := repo.IncreaseAccountBalance(ctx, credit.AccountID, credit.Amount)
		if err != nil {
			return fmt.Errorf("failed repo.DecreaseAccountBalance: %w", err)
		}

		userID, err := repo.SelectUserIDByAccountID(ctx, credit.AccountID)
		if err != nil {
			return fmt.Errorf("failed repo.SelectUserIDByAccountID: %w", err)
		}

		user, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: userID})

		msg := model.SimpleTelegramMessage{
			ChatID:  user.ChatId,
			Message: fmt.Sprintf("Кредит успешно получен! Сумма кредита: %d", credit.Amount),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.rbtmq.EnqueueTelegramMessage(ctx, msgBytes, nil)

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed processing in PaymentConsumeHandler: %w", err)
	}

	return nil
}
