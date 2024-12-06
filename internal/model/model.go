package model

import (
	"database/sql"
	"time"
)

type Payment struct {
	ID                string
	SenderAccountID   string
	ReceiverAccountID string
	PaymentID         string
	Type              string
	Amount            int64
	Currency          string
	Description       string
	Status            string
	CreatedAt         time.Time
	UpdatedAt         sql.NullTime
}

type Transaction struct {
	ID          string
	PaymentID   string
	AccountID   string
	Type        string
	Amount      int64
	Currency    string
	Description string
	Status      string
	CreatedAt   string
}
