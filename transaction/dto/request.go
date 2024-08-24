package dto

type AddTransactionRequest struct {
	Amount              float64 `json:"amount" validate:"required"`
	Type                string  `json:"type" validate:"required"`
	ParentTransactionId *uint64 `json:"parent_id"`
	TransactionId       string  `uri:"transactionId" validate:"required"`
}

type GetTransactionByIdRequest struct {
	TranscationId uint64 `uri:"transactionId" validate:"required"`
}

type GetTransactionByTypeRequest struct {
	TranscationType string `uri:"transactionType" validate:"required"`
}

