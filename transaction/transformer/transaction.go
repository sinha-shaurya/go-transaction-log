package transformer

import (
	"go-transaction-log/transaction/dto"
	"go-transaction-log/transaction/models"
)

func BuildGetTransactionResponse(transaction models.Transaction) dto.GetTransactionResponse{
	return dto.GetTransactionResponse{
		Id: transaction.ID,
		Amount: transaction.Amount,
		Type: transaction.Type,
		ParentTransactionId: transaction.ParentTransactionId,
		TransactionId: transaction.TransactionId,
	}
}

func BuildGetTransactionByTypeResponse(transactions []models.Transaction) []dto.GetTransactionResponse{
	responseList := make([]dto.GetTransactionResponse, 0)

	for _ , item := range transactions{
		responseList = append(responseList,BuildGetTransactionResponse(item))
	}

	return responseList
}