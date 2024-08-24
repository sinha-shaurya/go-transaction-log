package dto

type AddTransactionResponse struct {
	Status string `json:"status"`
}

type GetTransactionResponse struct {
	Id                  uint64  `json:"id"`
	Amount              float64 `json:"amount"`
	Type                string  `json:"type"`
	ParentTransactionId *uint64 `json:"parent_id"`
	TransactionId       string  `json:"transaction_id"`
}

type SumTransactionResponse struct {
	Sum float64 `json:"sum"`
}
