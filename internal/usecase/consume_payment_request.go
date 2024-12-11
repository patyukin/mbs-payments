package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) ConsumerPaymentRequest(ctx context.Context, record *kgo.Record) error {
	select {
	case <-ctx.Done():
		log.Error().Msgf("Context is done before processing in PaymentConsumeHandler: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	var message model.PaymentRequest
	if err := json.Unmarshal(record.Value, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message for topic '%s': %w", record.Topic, err)
	}

	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			payment, err := repo.GetPaymentByID(ctx, message.PaymentID)
			if err != nil {
				log.Error().Msgf("failed repo.GetPayment: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.GetPayment: %w", err)
			}

			err = repo.DecreaseAccountBalance(ctx, payment.SenderAccountID, payment.Amount)
			if err != nil {
				log.Error().Msgf("failed repo.UpdateAccountBalance: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.UpdateAccountBalance: %w", err)
			}

			_, err = repo.InsertTransaction(ctx, payment, payment.SenderAccountID, "DEBIT")
			if err != nil {
				log.Error().Msgf("failed repo.InsertTransaction: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.InsertTransaction: %w", err)
			}

			err = repo.IncreaseAccountBalance(ctx, payment.ReceiverAccountID, payment.Amount)
			if err != nil {
				log.Error().Msgf("failed repo.UpdateAccountBalance: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.UpdateAccountBalance: %w", err)
			}

			_, err = repo.InsertTransaction(ctx, payment, payment.ReceiverAccountID, "CREDIT")
			if err != nil {
				log.Error().Msgf("failed repo.InsertTransaction: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.InsertTransaction: %w", err)
			}

			err = repo.UpdatePaymentStatusByID(ctx, payment.ID, "COMPLETED")
			if err != nil {
				log.Error().Msgf("failed repo.UpdatePaymentStatus: %v", err)
				if err = repo.UpdatePaymentStatusByID(ctx, message.PaymentID, "FAILED"); err != nil {
					log.Error().Msgf("failed repo.UpdatePaymentStatus to status FAILED: %v", err)
				}

				return fmt.Errorf("failed repo.UpdatePaymentStatus: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	log.Info().Msg("consumer done")

	return nil
}
