package main

import (
	"context"
	"fmt"
	"github.com/patyukin/bs-payments/internal/config"
	"github.com/patyukin/bs-payments/internal/handler"
	"github.com/patyukin/bs-payments/internal/server"
	"github.com/patyukin/bs-payments/pkg/tracer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("unable to load config: %v", err)
	}

	srvAddress := fmt.Sprintf(":%d", cfg.HttpServer.Port)

	// register metrics
	err = metrics.RegisterMetrics()
	if err != nil {
		log.Fatal().Msgf("failed to register metrics: %v", err)
	}

	traceProvider, err := tracer.Init(fmt.Sprintf("%s/api/traces", cfg.TracerHost), "Api Gateway")
	if err != nil {
		log.Fatal().Msgf("failed init tracer, err: %v", err)
	}

	// auth service init
	authConn, err := grpc_client.NewGRPCClientServiceConn(cfg.GRPC.AuthServicePort)
	if err != nil {
		log.Fatal().Msgf("failed to connect to auth service: %v", err)
	}

	defer func(authConn *grpc.ClientConn) {
		if err = authConn.Close(); err != nil {
			log.Error().Msgf("failed to close auth service connection: %v", err)
		}
	}(authConn)

	authClient := authpb.NewAuthServiceClient(authConn)
	authUseCase := auth.New([]byte(cfg.JwtSecret), authClient)

	h := handler.New(authUseCase)
	r := server.InitRouterWithTrace(h, cfg, srvAddress)
	srv := server.New(r)

	errCh := make(chan error)

	go func() {
		log.Info().Msgf("starting server on %d", cfg.HttpServer.Port)
		if err = srv.Run(srvAddress, cfg); err != nil {
			log.Error().Msgf("failed starting server: %v", err)
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
			log.Info().Msgf("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msgf("Signal received")
		}
	}

	log.Info().Msgf("Shutting Down")

	if err = srv.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed server shutting down: %s", err.Error())
	}

	if err = traceProvider.Shutdown(ctx); err != nil {
		log.Error().Msgf("Error shutting down tracer provider: %v", err)
	}
}
