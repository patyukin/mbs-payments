package server

import (
	"context"
	paymentpb "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
)

type UseCase interface {
	CreateAccountUseCase(ctx context.Context, request *paymentpb.CreateAccountRequest) (*paymentpb.CreateAccountResponse, error)
	CreatePaymentUseCase(ctx context.Context, request *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error)
	ConfirmationPaymentUseCase(ctx context.Context, in *paymentpb.ConfirmationPaymentRequest) (*paymentpb.ConfirmationPaymentResponse, error)
}

type Server struct {
	paymentpb.UnimplementedPaymentServiceServer
	uc UseCase
}

func (s *Server) GetPayment(ctx context.Context, request *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) UpdatePaymentStatus(ctx context.Context, request *paymentpb.UpdatePaymentStatusRequest) (*paymentpb.UpdatePaymentStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetTransactionsByPayment(ctx context.Context, request *paymentpb.GetTransactionsByPaymentRequest) (*paymentpb.GetTransactionsByPaymentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func New(uc UseCase) *Server {
	return &Server{
		uc: uc,
	}
}
