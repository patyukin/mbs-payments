package server

import (
	"context"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

type UseCase interface {
	CreateAccountUseCase(ctx context.Context, request *paymentpb.CreateAccountRequest) (*paymentpb.CreateAccountResponse, error)
	CreatePaymentUseCase(ctx context.Context, request *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error)
	ConfirmationPaymentUseCase(ctx context.Context, in *paymentpb.ConfirmationPaymentRequest) (*paymentpb.ConfirmationPaymentResponse, error)
	GetPaymentUseCase(ctx context.Context, in *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error)
}

type Server struct {
	paymentpb.UnimplementedPaymentServiceServer
	uc UseCase
}

func New(uc UseCase) *Server {
	return &Server{
		uc: uc,
	}
}
