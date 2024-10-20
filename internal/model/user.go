package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type UserRole string

type SignUpData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Telegram string `json:"telegram"`
}

type SignUpV2Data struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CheckerPasswordData struct {
	UUID         uuid.UUID
	PasswordHash string
}

type SignInVerifyData struct {
	Code string `json:"code"`
}

type User struct {
	UUID      uuid.UUID      `json:"id,omitempty"`
	Login     string         `json:"login"`
	Name      sql.NullString `json:"name"`
	Surname   sql.NullString `json:"surname"`
	Role      UserRole       `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

type UserAuthInfo struct {
	UserUUID uuid.UUID `json:"user_id"`
	Role     UserRole  `json:"role"`
}
