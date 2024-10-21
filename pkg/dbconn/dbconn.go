package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func New(ctx context.Context, cfg PostgreSQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	log.Info().Msg("connected to db")

	err = dbConn.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	log.Info().Msg("pinged postgresql db")

	return dbConn, nil
}
