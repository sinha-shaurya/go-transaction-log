package repository

import (
	"errors"
	"go-transaction-log/transaction/models"
	"go-transaction-log/utils"

	"gorm.io/gorm"
)

func CreateTransaction(db *gorm.DB, transactionModel *models.Transaction) error {
	return db.Create(transactionModel).Error
}

func GetTransactionByFilter(db *gorm.DB, filter TransactionLogFilter) ([]models.Transaction, error) {
	if utils.CheckIfStructFieldsAreNil(&filter) {
		return nil, errors.New("no DB Filters found. Aborting to prevent entire table scan")
	}

	var transactions []models.Transaction
	result := db.Scopes(withId(filter.Id),
		withTransactionId(filter.TransactionId),
		withType(filter.Type)).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return transactions, nil
}

func GetTransactionByIds(db *gorm.DB, ids []uint64) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := db.Scopes(withIdIn(ids)).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return transactions, nil
}
