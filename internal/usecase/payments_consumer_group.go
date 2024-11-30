package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) PaymentsConsumerGroup(ctx context.Context, record *kgo.Record) error {
	select {
	case <-ctx.Done():
		log.Error().Msgf("Context is done before processing in PaymentsConsumerGroup: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	switch record.Topic {
	case kafka.TransactionsTopic:
		err := u.UpdateTransactionsStatus(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.UpdateTransactionsStatus: %w", err)
		}
	case kafka.PaymentRequestTopic:
		err := u.ConsumerPaymentRequest(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	case kafka.CreditCreatedTopic:
		err := u.ConsumerCreditCreated(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	case kafka.CreditPaymentsTopic:
		err := u.ConsumerCreditPayments(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	default:
		return fmt.Errorf("failed to unmarshal message for topic '%s'", record.Topic)
	}

	return nil
}
