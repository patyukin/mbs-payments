package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/bs-payments/internal/db"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (u *UseCase) GetPaymentUseCase(ctx context.Context, in *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {
		payment, err := repo.GetPaymentByID(ctx, in.GetPaymentId())
		if err != nil {
			return fmt.Errorf("failed u.registry.GetRepo().GetPaymentByID: %w", err)
		}

		userID, err := repo.SelectUserIDByAccountID(ctx, payment.SenderAccountID)
		if err != nil {
			return fmt.Errorf("failed u.registry.GetRepo().SelectUserIDByAccountID: %w", err)
		}

		if userID != in.GetUserId() {
			return fmt.Errorf("failed u.registry.GetRepo().SelectUserIDByAccountID: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed u.registry.ReadCommitted: %w", err)
	}

	return nil, nil
}
