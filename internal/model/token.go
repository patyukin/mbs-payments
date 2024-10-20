package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token     uuid.UUID    `json:"token"`
	UserUUID  uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
