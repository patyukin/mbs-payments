package cacher

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type Cacher struct {
	client *redis.Client
}

func New(ctx context.Context, dsn string) (*Cacher, error) {
	c := redis.NewClient(&redis.Options{Addr: dsn})

	err := c.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Info().Msg("connected to redis")
	return &Cacher{client: c}, nil
}

func (c *Cacher) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}

	return nil
}

func (c *Cacher) SetPaymentConfirmationCode(ctx context.Context, userID, paymentID, code string) error {
	return c.client.Set(ctx, fmt.Sprintf("u:%s:pc:%s", userID, paymentID), code, 50*time.Minute).Err()
}

func (c *Cacher) GetPaymentConfirmationCode(ctx context.Context, userID, code string) (string, error) {
	pattern := fmt.Sprintf("u:%s:pc:*", userID)
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get keys: %w", err)
	}

	var storedCode string
	for _, key := range keys {
		storedCode, err = c.client.Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			continue
		}

		if err != nil {
			return "", fmt.Errorf("failed to get value for key %s: %w", key, err)
		}

		if storedCode == code {
			return key[len(fmt.Sprintf("u:%s:pc:", userID)):], nil
		}
	}

	return "", fmt.Errorf("payment confirmation code not found for userID: %s and code: %s", userID, code)
}

func (c *Cacher) DeletePaymentConfirmationCode(ctx context.Context, userID, paymentID string) error {
	err := c.client.Del(ctx, fmt.Sprintf("u:%s:pc:%s", userID, paymentID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete payment confirmation code: %w", err)
	}

	return nil
}
