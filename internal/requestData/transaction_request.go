package requestData

type TransactionReq struct {
	AccountId      int64   `json:"accountId"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}
