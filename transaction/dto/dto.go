package dto

type AddTransactionRequestDto struct {
	Amount              float64
	Type                string
	ParentTransactionId *uint64
	TransactionId       string
}
