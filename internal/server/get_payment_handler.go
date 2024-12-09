package server

import (
	"context"
	"fmt"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (s *Server) GetPayment(ctx context.Context, in *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	payment, err := s.uc.GetPaymentUseCase(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed uc.CreatePaymentUseCase: %w", err)
	}

	return payment, nil
}
