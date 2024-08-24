package models

import (
	"gorm.io/datatypes"
)

type Transaction struct {
	ID                  uint64  `gorm:"primaryKey;autoincrement"`
	TransactionId       string  `gorm:"type:varchar(50);not null;unique"`
	Amount              float64 `gorm:"type:double"`
	Type                string  `gorm:"type:varchar(50);index:idx_type"`
	ParentTransactionId *uint64 `gorm:"type:varchar(50);default:NULL"`
	Meta                datatypes.JSONType[TransactionMeta]
}

type TransactionMeta struct {
	ParentIds []uint64 `json:"parent_ids,omitempty"`
}

func (t *Transaction) GetMeta() TransactionMeta {
	return t.Meta.Data()
}

func (t *Transaction) SetMeta(meta TransactionMeta) {
	t.Meta = datatypes.NewJSONType[TransactionMeta](meta)
}
