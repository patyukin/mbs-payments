package server

import (
	"context"
	"fmt"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/rs/zerolog/log"
)

func (s *Server) CreatePayment(ctx context.Context, in *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	log.Debug().Msgf("in handler: %v", in)
	result, err := s.uc.CreatePaymentUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.CreatePaymentUseCase: %w", err)
	}

	return result, nil
}
