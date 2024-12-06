package usecase

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) PaymentsConsumerGroup(ctx context.Context, record *kgo.Record) error {
	log.Debug().Msgf("Received to topic: %s, record: %v", record.Topic, string(record.Value))

	switch record.Topic {
	case kafka.TransactionReportSolutionTopic:
		err := u.ConsumeTransactionReportSolution(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumeTransactionReportSolution: %w", err)
		}
	case kafka.PaymentRequestTopic:
		err := u.ConsumerPaymentRequest(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	case kafka.CreditCreatedTopic:
		err := u.ConsumeCreditCreated(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	case kafka.CreditPaymentsTopic:
		err := u.ConsumeCreditPayments(ctx, record)
		if err != nil {
			return fmt.Errorf("failed u.ConsumerPaymentRequest: %w", err)
		}
	default:
		return fmt.Errorf("failed to unmarshal message for topic '%s'", record.Topic)
	}

	return nil
}
