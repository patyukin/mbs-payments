package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (u *UseCase) ConfirmationPaymentUseCase(ctx context.Context, in *paymentpb.ConfirmationPaymentRequest) (*paymentpb.ConfirmationPaymentResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		paymentID, err := u.chr.GetPaymentConfirmationCode(ctx, in.UserId, in.Code)
		if err != nil {
			return fmt.Errorf("failed u.chr.GetPaymentConfirmationCode: %w", err)
		}

		payment, err := repo.GetPaymentByID(ctx, paymentID)
		if err != nil {
			return fmt.Errorf("failed repo.GetPayment: %w", err)
		}

		balance, err := repo.ExistsPositiveAccountBalance(ctx, payment.SenderAccountID)
		if err != nil {
			return fmt.Errorf("failed repo.InsertPayment: %w", err)
		}

		if balance < payment.Amount {
			return fmt.Errorf("not enough balance")
		}

		err = repo.UpdatePaymentStatusByID(ctx, paymentID, "PENDING")
		if err != nil {
			return fmt.Errorf("failed repo.UpdatePaymentStatus: %w", err)
		}

		paymentRequest := model.PaymentRequest{PaymentID: paymentID}
		paymentRequestBytes, err := json.Marshal(paymentRequest)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		err = u.kfk.PublishPaymentReport(ctx, paymentRequestBytes)
		if err != nil {
			return fmt.Errorf("failed u.rbtmq.PushPaymentStatusChange: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return &paymentpb.ConfirmationPaymentResponse{Message: "Платеж успешно подтвержден"}, nil
}
