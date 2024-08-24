package repository

import (
	"gorm.io/gorm"
)

func withId(id *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id != nil {
			return db.Where("id = ?", *id)
		}

		return db
	}
}

func withIdIn(ids []uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) != 0 {
			return db.Where("id IN ?", ids)
		}

		return db
	}
}


func withTransactionId(transactionId *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if transactionId != nil {
			return db.Where("transaction_id = ?", transactionId)
		}

		return db
	}
}

func withType(transactionType *string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if transactionType != nil {
			return db.Where("type = ?", transactionType)
		}

		return db
	}
}
