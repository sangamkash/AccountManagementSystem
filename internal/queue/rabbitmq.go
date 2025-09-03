package queue

// TransactionMessage defines the payload
type TransactionMessage struct {
	IdempotencyKey string  `json:"idempotency_key"`
	AccountID      int64   `json:"account_id"`
	Type           string  `json:"type"`
	Amount         float64 `json:"amount"`
}
