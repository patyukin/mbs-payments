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

func (u *UseCase) ConsumeCreditPayments(ctx context.Context, record *kgo.Record) error {
	var creditPaymentsSolution []model.CreditPaymentSolution

	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			var payments []model.CreditPayment
			if err := json.Unmarshal(record.Value, &payments); err != nil {
				return fmt.Errorf("failed to unmarshal messages: %w", err)
			}

			for _, payment := range payments {
				status := "COMPLETED"
				if err := repo.DecreaseAccountBalance(ctx, payment.AccountID, payment.Amount); err != nil {
					log.Error().Msgf("failed repo.UpdateAccountBalance: %v", err)
					status = "FAILED"
				}

				creditPaymentsSolution = append(
					creditPaymentsSolution, model.CreditPaymentSolution{
						PaymentScheduleID: payment.PaymentScheduleID, Status: status,
					},
				)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed processing in PaymentConsumeHandler: %w", err)
	}

	bytes, err := json.Marshal(creditPaymentsSolution)
	if err != nil {
		return fmt.Errorf("failed to marshal messages: %w", err)
	}

	if err = u.kfk.PublishCreditPaymentsSolution(ctx, bytes); err != nil {
		return fmt.Errorf("failed u.kafkaProducer.PublishCreditPaymentSolution: %w", err)
	}

	return nil
}
