package service

import (
	errors "errors"
	"go-transaction-log/transaction/dto"
	"go-transaction-log/transaction/models"
	"go-transaction-log/transaction/repository"
	httperrors "go-transaction-log/utils/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AddTransaction(ctx *gin.Context, db *gorm.DB, requestDto dto.AddTransactionRequestDto) *httperrors.HttpError {

	logger, err := zap.NewProduction()
	if err != nil {
		return &httperrors.HttpError{
			ErrorCode: int(http.StatusInternalServerError),
			Error:     err,
		}
	}

	ancestorIds, apiErr := validateParentTransactionId(db, requestDto.ParentTransactionId)
	if apiErr != nil {
		return apiErr
	}
	meta := buildMetaForTransaction(requestDto.ParentTransactionId, ancestorIds)

	transactionModel := buildTransactionFromRequest(requestDto)
	transactionModel.SetMeta(meta)
	logger.Info("Created model")

	createErr := createTransaction(db, transactionModel)

	if createErr != nil {
		return &httperrors.HttpError{
			ErrorCode: int(http.StatusInternalServerError),
			Error:     createErr,
		}
	}

	return nil
}

func validateParentTransactionId(db *gorm.DB, parentTransactionId *uint64) ([]uint64, *httperrors.HttpError) {
	if parentTransactionId == nil {
		return nil, nil
	}

	parentTransactionMeta, apiErr := getParentTransactionMeta(db, parentTransactionId)
	if apiErr != nil {
		return nil, apiErr
	}
	return parentTransactionMeta.ParentIds, nil
}

func buildTransactionFromRequest(request dto.AddTransactionRequestDto) *models.Transaction {

	t := models.Transaction{
		TransactionId: request.TransactionId,
		Amount:        request.Amount,
		Type:          request.Type,
	}

	if request.ParentTransactionId != nil {
		t.ParentTransactionId = request.ParentTransactionId
	}
	return &t
}

func createTransaction(db *gorm.DB, transactionModel *models.Transaction) error {
	createErr := repository.CreateTransaction(db, transactionModel)
	if createErr != nil {
		return createErr
	}

	return nil
}

func getParentTransactionMeta(db *gorm.DB, transactionId *uint64) (*models.TransactionMeta, *httperrors.HttpError) {
	filter := repository.TransactionLogFilter{
		Id: transactionId,
	}

	parentTransaction, err := repository.GetTransactionByFilter(db, filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &httperrors.HttpError{
				ErrorCode: http.StatusNotFound,
				Error:     errors.New("No record found with given parent transaction id"),
			}
		}
	}

	meta := parentTransaction[0].GetMeta()
	return &meta, nil
}

func buildMetaForTransaction(parentTransactionId *uint64, ancestorIds []uint64) models.TransactionMeta {
	parentTransactionIds := make([]uint64, 0)
	if parentTransactionId == nil {
		return models.TransactionMeta{}
	}
	if len(ancestorIds) != 0 {
		parentTransactionIds = append(parentTransactionIds, ancestorIds...)
	}

	return models.TransactionMeta{ParentIds: parentTransactionIds}
}

func GetTransactionById(db *gorm.DB, id uint64) (*models.Transaction, *httperrors.HttpError) {
	transaction, dbErr := repository.GetTransactionByFilter(db, repository.TransactionLogFilter{
		Id: &id,
	})

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			return nil, &httperrors.HttpError{
				ErrorCode: http.StatusNotFound,
				Error:     errors.New("No record found"),
			}
		}
		return nil, &httperrors.HttpError{
			ErrorCode: http.StatusInternalServerError,
			Error:     dbErr,
		}
	}

	return &(transaction[0]), nil
}

func GetTransactionByType(db *gorm.DB, transactionType string) ([]models.Transaction, *httperrors.HttpError) {
	transactions, dbErr := repository.GetTransactionByFilter(db, repository.TransactionLogFilter{
		Type: &transactionType,
	})

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			return nil, &httperrors.HttpError{
				ErrorCode: http.StatusNotFound,
				Error:     errors.New("No record found"),
			}
		}
		return nil, &httperrors.HttpError{
			ErrorCode: http.StatusInternalServerError,
			Error:     dbErr,
		}
	}

	return transactions, nil
}

func SumTransactions(db *gorm.DB, id uint64) (float64, *httperrors.HttpError) {
	// find transaction by id
	transcations, dbErr := repository.GetTransactionByFilter(db, repository.TransactionLogFilter{
		Id: &id,
	})
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			return 0, &httperrors.HttpError{
				ErrorCode: http.StatusNotFound,
				Error:     errors.New("No record found"),
			}
		}
		return 0, &httperrors.HttpError{
			ErrorCode: http.StatusInternalServerError,
			Error:     dbErr,
		}
	}

	parentIds := transcations[0].GetMeta().ParentIds
	if len(parentIds) == 0 {
		return transcations[0].Amount, nil
	}

	// get all parent ids
	ancestorAmounts, err := computeAmountForAncestorTransactions(db, parentIds)
	if err !=nil{
		return 0.0, err
	}

	return transcations[0].Amount + ancestorAmounts , nil
}

func computeAmountForAncestorTransactions(db *gorm.DB, ancestorIds []uint64) (float64, *httperrors.HttpError){
	ancestorTransactions, dbErr := repository.GetTransactionByIds(db, ancestorIds)
	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			return 0.0, &httperrors.HttpError{
				ErrorCode: http.StatusNotFound,
				Error:     errors.New("Invalid parent ids in transaction"),
			}
		}
		return 0.0, &httperrors.HttpError{
			ErrorCode: http.StatusInternalServerError,
			Error:     dbErr,
		}
	}

	amount := 0.0
	for _ , transaction := range ancestorTransactions{
		amount += transaction.Amount
	}

	return amount, nil
}