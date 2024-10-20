package db

import (
	"auth-telegram/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (r *Repository) InsertIntoUser(ctx context.Context, in model.SignUpData) (uuid.UUID, error) {
	currentTime := time.Now().UTC()
	query := `INSERT INTO users (login, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.Login, in.Password, currentTime)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert user: %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}

func (r *Repository) SelectUserByLogin(ctx context.Context, login string) (model.CheckerPasswordData, error) {
	query := `SELECT id, password_hash FROM users WHERE login = $1`
	row := r.db.QueryRowContext(ctx, query, login)
	if row.Err() != nil {
		return model.CheckerPasswordData{}, fmt.Errorf("failed to select user: %w", row.Err())
	}

	var user model.CheckerPasswordData
	err := row.Scan(&user.UUID, &user.PasswordHash)
	if err != nil {
		return model.CheckerPasswordData{}, fmt.Errorf("failed to select user: %w", err)
	}

	return user, nil
}

func (r *Repository) SelectUserByUUID(ctx context.Context, userUUID string) (model.User, error) {
	query := `SELECT id, login, name, surname, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID)
	if row.Err() != nil {
		return model.User{}, fmt.Errorf("failed to select user: %w", row.Err())
	}

	var user model.User
	err := row.Scan(&user.UUID, &user.Login, &user.Name, &user.Surname, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to select user: %w", err)
	}

	return user, nil
}

func (r *Repository) InsertToken(ctx context.Context, userUUID uuid.UUID) (uuid.UUID, error) {
	currentTime := time.Now().UTC()
	expiresAt := currentTime.Add(24 * 30 * time.Hour)
	query := `INSERT INTO tokens (user_id, expires_at, created_at) VALUES ($1, $2, $3) RETURNING token`
	row := r.db.QueryRowContext(ctx, query, userUUID.String(), expiresAt, currentTime)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert token: %w", row.Err())
	}

	var token uuid.UUID
	err := row.Scan(&token)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert token: %w", err)
	}

	return token, nil
}

func (r *Repository) SelectFromTelegramUsersByUser(ctx context.Context, userUUID uuid.UUID) (model.TelegramUser, error) {
	query := `SELECT id, user_id, tg_username,  tg_user_id, tg_chat_id, created_at, updated_at FROM telegram_users WHERE user_id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID.String())
	if row.Err() != nil {
		return model.TelegramUser{}, fmt.Errorf("failed to select telegram_users: %w", row.Err())
	}

	var tgUser model.TelegramUser
	err := row.Scan(&tgUser.UUID, &tgUser.UserUUID, &tgUser.TgUsername, &tgUser.TgUserID, &tgUser.TgChatID, &tgUser.CreatedAt, &tgUser.UpdatedAt)
	if err != nil {
		return model.TelegramUser{}, fmt.Errorf("failed to select telegram_users: %w", err)
	}

	return tgUser, nil

}

func (r *Repository) InsertIntoTelegramUsers(ctx context.Context, telegramLogin string, userUUID uuid.UUID) (uuid.UUID, error) {
	currentTime := time.Now().UTC()
	query := `INSERT INTO telegram_users (user_id, tg_username, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, userUUID.String(), telegramLogin, currentTime)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert in telegram_users: %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert in telegram_users: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateTelegramUserAfterSignUp(ctx context.Context, userUUID uuid.UUID, chatID, userID int64) error {
	currentTime := time.Now().UTC()
	query := `UPDATE telegram_users SET tg_user_id = $1, tg_chat_id = $2, updated_at = $3 WHERE user_id = $4`
	_, err := r.db.ExecContext(ctx, query, userID, chatID, currentTime, userUUID.String())
	if err != nil {
		return fmt.Errorf("failed to insert in telegram_users: %w", err)
	}

	return nil
}

func (r *Repository) SelectUserAuthInfoByUUID(ctx context.Context, userUUID string) (model.UserAuthInfo, error) {
	query := `SELECT id, role FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userUUID)
	if row.Err() != nil {
		return model.UserAuthInfo{}, fmt.Errorf("failed to select user: %w", row.Err())
	}

	var id uuid.UUID
	var role string
	err := row.Scan(&id, &role)
	if err != nil {
		return model.UserAuthInfo{}, fmt.Errorf("failed to select user: %w", err)
	}

	return model.UserAuthInfo{UserUUID: id, Role: model.UserRole(role)}, nil
}

func (r *Repository) InsertIntoUserV2(ctx context.Context, in model.SignUpV2Data) (uuid.UUID, error) {
	currentTime := time.Now().UTC()
	query := `INSERT INTO users (login, password_hash, created_at) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, in.Login, in.Password, currentTime)
	if row.Err() != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert user: %w", row.Err())
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return id, nil
}
