package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) CreateAccountUseCase(ctx context.Context, in *paymentpb.CreateAccountRequest) (*paymentpb.CreateAccountResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		accountID, err := repo.InsertAccount(ctx, in)
		if err != nil {
			return fmt.Errorf("failed repo.InsertAccount: %w", err)
		}

		userInfo, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: in.UserId})
		if err != nil {
			return fmt.Errorf("failed u.authClient.GetUserInfo: %w", err)
		}

		log.Debug().Msgf("accountID: %s", accountID)

		msg := model.SimpleTelegramMessage{
			Message: fmt.Sprintf("создан счет: %s", accountID),
			ChatID:  userInfo.ChatId,
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.rbtmq.EnqueueTelegramMessage(ctx, msgBytes, nil)
		if err != nil {
			return fmt.Errorf("failed u.rbtmq.PublishAccountCreation: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &paymentpb.CreateAccountResponse{Message: "Счет успешно добавлен"}, nil
}
