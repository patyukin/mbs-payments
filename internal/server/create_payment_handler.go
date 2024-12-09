package server

import (
	"context"
	"fmt"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (s *Server) CreatePayment(ctx context.Context, in *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	result, err := s.uc.CreatePaymentUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.CreatePaymentUseCase: %w", err)
	}

	return result, nil
}
