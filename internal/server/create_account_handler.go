package server

import (
	"context"
	"fmt"
	"github.com/patyukin/mbs-pkg/pkg/errs"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

func (s *Server) CreateAccount(ctx context.Context, request *paymentpb.CreateAccountRequest) (*paymentpb.CreateAccountResponse, error) {
	result, err := s.uc.CreateAccountUseCase(ctx, request)
	if err != nil {
		return &paymentpb.CreateAccountResponse{
			Error: errs.ToErrorResponse(fmt.Errorf("failed s.uc.CreateAccountUseCase: %w", err)),
		}, nil
	}

	return result, nil
}
