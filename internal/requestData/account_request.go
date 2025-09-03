package requestData

type CreateAccountReq struct {
	Username      string  `json:"username"`
	InitialAmount float64 `json:"initial_amount"`
}
