package server

import (
	"context"
	"fmt"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) ConfirmationPayment(ctx context.Context, in *paymentpb.ConfirmationPaymentRequest) (*paymentpb.ConfirmationPaymentResponse, error) {
	log.Debug().Msgf("in handler: %v", in)
	result, err := s.uc.ConfirmationPaymentUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed s.uc.ConfirmationPaymentUseCase: %w", err)
	}

	return result, nil
}
