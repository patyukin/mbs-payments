package usecase

import (
	"context"
	"github.com/patyukin/bs-payments/internal/db"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type Cacher interface {
	SetPaymentConfirmationCode(ctx context.Context, userID, paymentID, code string) error
	GetPaymentConfirmationCode(ctx context.Context, userID, code string) (string, error)
	DeletePaymentConfirmationCode(ctx context.Context, userID, code string) error
}

type RabbitMQProducer interface {
	EnqueueTelegramMessage(ctx context.Context, body []byte, headers amqp.Table) error
}

type AuthClient interface {
	GetBriefUserByID(ctx context.Context, in *authpb.GetBriefUserByIDRequest, opts ...grpc.CallOption) (*authpb.GetBriefUserByIDResponse, error)
}

type KafkaProducer interface {
	PublishTransactionReport(ctx context.Context, value []byte) error
	PublishPaymentRequest(ctx context.Context, value []byte) error
	PublishCreditPaymentsSolution(ctx context.Context, value []byte) error
}

type UseCase struct {
	registry   *db.Registry
	rbtmq      RabbitMQProducer
	chr        Cacher
	kfk        KafkaProducer
	authClient AuthClient
}

func New(registry *db.Registry, rbtmq RabbitMQProducer, chr Cacher, kfk KafkaProducer, authClient AuthClient) *UseCase {
	return &UseCase{
		registry:   registry,
		rbtmq:      rbtmq,
		chr:        chr,
		kfk:        kfk,
		authClient: authClient,
	}
}
