package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/patyukin/bs-payments/internal/cacher"
	"github.com/patyukin/bs-payments/internal/config"
	"github.com/patyukin/bs-payments/internal/cronjob"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/bs-payments/internal/metrics"
	"github.com/patyukin/bs-payments/internal/server"
	"github.com/patyukin/bs-payments/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/dbconn"
	"github.com/patyukin/mbs-pkg/pkg/grpc_client"
	"github.com/patyukin/mbs-pkg/pkg/grpc_server"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/patyukin/mbs-pkg/pkg/migrator"
	"github.com/patyukin/mbs-pkg/pkg/mux_server"
	authpb "github.com/patyukin/mbs-pkg/pkg/proto/auth_v1"
	desc "github.com/patyukin/mbs-pkg/pkg/proto/payment_v1"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/patyukin/mbs-pkg/pkg/tracing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const ServiceName = "PaymentService"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	if err = metrics.Init(); err != nil {
		log.Fatal().Msgf("failed to init metrics: %v", err)
	}

	_, closer, err := tracing.InitJaeger(fmt.Sprintf("jaeger:6831"), ServiceName)
	if err != nil {
		log.Fatal().Msgf("failed to initialize tracer: %v", err)
	}

	defer closer()

	log.Info().Msg("Jaeger connected")

	log.Info().Msg("Opentracing connected")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	dbConn, err := dbconn.New(ctx, cfg.PostgreSQLDSN)
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	if err = migrator.UpMigrations(ctx, dbConn); err != nil {
		log.Fatal().Msgf("failed to up migrations: %v", err)
	}

	rbt, err := rabbitmq.New(cfg.RabbitMQURL, rabbitmq.Exchange)
	if err != nil {
		log.Fatal().Msgf("failed to create rabbit producer: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.TelegramMessageQueue,
		[]string{rabbitmq.TelegramMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind TelegramMessageQueue to exchange with - TelegramMessageRouteKey: %v", err)
	}

	chr, err := cacher.New(ctx, cfg.RedisDSN)
	if err != nil {
		log.Fatal().Msgf("failed to create redis cacher: %v", err)
	}

	kfk, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka consumer, err: %v", err)
	}

	kfkConsumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, cfg.Kafka.Topics)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka consumer, err: %v", err)
	}

	defer kfkConsumer.Close()

	// auth service init
	authConn, err := grpc_client.NewGRPCClientServiceConn(cfg.GRPC.AuthService)
	if err != nil {
		log.Fatal().Msgf("failed to connect to auth service: %v", err)
	}

	defer func(authConn *grpc.ClientConn) {
		if err = authConn.Close(); err != nil {
			log.Error().Msgf("failed to close auth service connection: %v", err)
		}
	}(authConn)

	authClient := authpb.NewAuthServiceClient(authConn)

	registry := db.New(dbConn)
	uc := usecase.New(registry, rbt, chr, kfk, authClient)
	srv := server.New(uc)

	// grpc server
	s := grpc_server.NewGRPCServer()
	reflection.Register(s)
	desc.RegisterPaymentServiceServer(s, srv)
	grpcPrometheus.Register(s)

	// mux server
	muxServer := mux_server.New()

	errCh := make(chan error)

	// cron job
	cj := cronjob.New(uc)
	go func() {
		if err = cj.Run(ctx); err != nil {
			log.Error().Msgf("failed adding cron job, err: %v", err)
			errCh <- err
		}
	}()

	// run payments consumer
	go func() {
		if err = kfkConsumer.ProcessMessages(ctx, uc.PaymentsConsumerGroup); err != nil {
			log.Error().Msgf("failed to process messages: %v", err)
			errCh <- err
		}
	}()

	// GRPC server
	go func() {
		log.Info().Msgf("GRPC started on :%d", cfg.GRPCServer.Port)
		if err = s.Serve(lis); err != nil {
			log.Error().Msgf("failed to serve: %v", err)
			errCh <- err
		}
	}()

	// metrics + pprof server
	go func() {
		log.Info().Msgf("Prometheus metrics exposed on :%d/metrics", cfg.HttpServer.Port)
		if err = muxServer.Run(cfg.HttpServer.Port); err != nil {
			log.Error().Msgf("Failed to serve Prometheus metrics: %v", err)
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msg("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msg("Signal received")
		}
	}

	log.Info().Msg("Shutting Down")

	// stop server
	s.GracefulStop()

	if err = muxServer.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed to shutdown http server: %s", err.Error())
	}

	if err = dbConn.Close(); err != nil {
		log.Error().Msgf("failed db connection close: %s", err.Error())
	}

	if err = chr.Close(); err != nil {
		log.Error().Msgf("failed redis connection close: %s", err.Error())
	}
}
