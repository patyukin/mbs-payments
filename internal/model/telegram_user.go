package model

import (
	"github.com/google/uuid"
	"time"
)

type TelegramUser struct {
	UUID       uuid.UUID `json:"id"`
	UserUUID   uuid.UUID `json:"user_id"`
	TgUsername string    `json:"tg_username"`
	TgUserID   int64     `json:"tg_user_id"`
	TgChatID   int64     `json:"tg_chat_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
