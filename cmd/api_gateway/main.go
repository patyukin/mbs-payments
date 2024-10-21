package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/patyukin/bs-payments/internal/config"
	"github.com/patyukin/bs-payments/internal/db"
	"github.com/patyukin/bs-payments/internal/migrator"
	"github.com/patyukin/bs-payments/internal/usecase"
	"github.com/patyukin/bs-payments/pkg/dbconn"
	"github.com/patyukin/bs-payments/pkg/otely"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCServer.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	dbConn, err := dbconn.New(ctx, dbconn.PostgreSQLConfig(cfg.PostgreSQL))
	if err != nil {
		log.Fatal().Msgf("failed to connect to db: %v", err)
	}

	if err = migrator.UpMigrations(ctx, dbConn); err != nil {
		log.Fatal().Msgf("failed to up migrations: %v", err)
	}

	registry := db.New(dbConn)
	uc := usecase.New(registry)
	desc.RegisterAuthServiceServer(s, srv)

	// Создаем экземпляр MetricsWrapper
	otelyShutdown, err := otely.SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to initialize metrics wrapper: %v", err)
	}

	wg := &sync.WaitGroup{}

	// GRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Info().Msgf("GRPC started on :%d", cfg.GRPCServer.Port)
		if err = s.Serve(lis); err != nil {
			log.Fatal().Msgf("failed to serve: %v", err)
		}
	}()

	// metrics server
	wg.Add(1)
	go func() {
		defer wg.Done()

		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msgf("Prometheus metrics exposed on :%d/metrics", cfg.HttpServer.Port)
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpServer.Port), nil); err != nil {
			log.Fatal().Msgf("Failed to serve Prometheus metrics: %v", err)
		}
	}()

	wg.Wait()

	err = errors.Join(err, otelyShutdown(ctx))
}
