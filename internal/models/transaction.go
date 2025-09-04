package models

import (
	"time"
)

type Transaction struct {
	ID             int64     `json:"id" db:"id"`
	AccountID      int64     `json:"account_id" db:"account_id"`
	Type           string    `json:"type" db:"type"`
	Amount         float64   `json:"amount" db:"amount"`
	IdempotencyKey string    `json:"idempotency_key" db:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type TransactionMessage struct {
	IdempotencyKey string  `json:"idempotency_key"`
	AccountID      int64   `json:"account_id"`
	Type           string  `json:"type"`
	Amount         float64 `json:"amount"`
}
