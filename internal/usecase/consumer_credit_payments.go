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

func (u *UseCase) ConsumerCreditPayments(ctx context.Context, record *kgo.Record) error {
	err := u.registry.ReadCommitted(ctx, func(ctx context.Context, repo *db.Repository) error {

		var payments []model.CreditPayment
		if err := json.Unmarshal(record.Value, &payments); err != nil {
			return fmt.Errorf("failed to unmarshal messages: %w", err)
		}

		creditPaymentsSolution := make([]model.CreditPaymentSolution, 0, len(payments))

		for _, payment := range payments {
			status := "COMPLETED"
			err := repo.DecreaseAccountBalance(ctx, payment.AccountID, payment.Amount)
			if err != nil {
				log.Error().Msgf("failed repo.UpdateAccountBalance: %v", err)
				status = "FAILED"
			}

			creditPaymentsSolution = append(creditPaymentsSolution, model.CreditPaymentSolution{
				PaymentScheduleID: payment.PaymentScheduleID, Status: status,
			})
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed processing in PaymentConsumeHandler: %w", err)
	}

	return nil
}
