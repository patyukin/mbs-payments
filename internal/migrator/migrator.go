package migrator

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func UpMigrations(ctx context.Context, db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "./migrations"); err != nil {
		return fmt.Errorf("failed to up migrations: %w", err)
	}

	log.Info().Msg("up migrations")

	return nil
}
