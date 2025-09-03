package requestData

type createAccountReq struct {
	Username      string  `json:"username"`
	InitialAmount float64 `json:"initial_amount"`
}
