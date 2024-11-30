package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/rs/zerolog/log"
)

func (u *UseCase) CreatePaymentUseCase(ctx context.Context, in *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	log.Debug().Msgf("in: %v", in)

	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		balance, err := repo.ExistsPositiveAccountBalance(ctx, in.SenderAccountId)
		if err != nil {
			return fmt.Errorf("failed repo.InsertPayment: %w", err)
		}

		if balance < in.Amount {
			return fmt.Errorf("not enough balance")
		}

		paymentID, err := repo.InsertPayment(ctx, in)
		if err != nil {
			return fmt.Errorf("failed repo.InsertPayment: %w", err)
		}

		userInfo, err := u.authClient.GetBriefUserByID(ctx, &authpb.GetBriefUserByIDRequest{UserId: in.UserId})
		if err != nil {
			return fmt.Errorf("failed u.authClient.GetUserInfo: %w", err)
		}

		log.Debug().Msgf("accountID: %s", paymentID)

		code, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("failed to generate UUID: %w", err)
		}

		err = u.chr.SetPaymentConfirmationCode(ctx, paymentID, code.String(), in.SenderAccountId)
		if err != nil {
			return fmt.Errorf("failed u.chr.SetPaymentConfirmationCode: %w", err)
		}

		msg := model.SimpleTelegramMessage{
			Message: fmt.Sprintf("подтвердите платеж с кодом в течение 1 минуты: %s", code.String()),
			ChatID:  userInfo.ChatId,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.rbtmq.PublishPaymentExecutionInitiate(ctx, msgBytes, nil)
		if err != nil {
			return fmt.Errorf("failed u.rbtmq.PublishPaymentExecutionInitiate: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &paymentpb.CreatePaymentResponse{}, nil
}
