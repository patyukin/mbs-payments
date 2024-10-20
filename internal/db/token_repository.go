package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Repository) CleanTokens(ctx context.Context) error {
	currentTime := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, "DELETE FROM tokens WHERE expires_at < $1", currentTime)
	if err != nil {
		return fmt.Errorf("failed cleaning tokens: %w", err)
	}

	return nil
}

func (r *Repository) GetUserUUIDByRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	query := `SELECT user_id FROM tokens WHERE token = $1`
	row := r.db.QueryRowContext(ctx, query, refreshToken)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed to select token: %w", row.Err())
	}

	var userUUID uuid.UUID
	err := row.Scan(&userUUID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to select token: %w", err)
	}

	return userUUID, nil
}

func (r *Repository) DeleteToken(ctx context.Context, refreshToken string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tokens WHERE token = $1", refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (r *Repository) SelectServiceTokenByName(ctx context.Context, name string) (string, error) {
	query := `SELECT secret FROM services WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)
	if row.Err() != nil {
		return "", fmt.Errorf("failed to select token: %w", row.Err())
	}

	var hash string
	err := row.Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("failed to select token: %w", err)
	}

	return hash, nil
}
