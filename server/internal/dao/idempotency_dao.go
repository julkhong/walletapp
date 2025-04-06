package dao

import (
	"errors"

	"gorm.io/gorm"
)

func (dao *WalletDao) CheckIdempotencyKey(key, method, path string) (*IdempotencyRecord, bool) {
	var record IdempotencyRecord
	result := dao.db.
		Table("idempotency_keys").
		Where("key = ? AND method = ? AND path = ?", key, method, path).
		First(&record)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false
	}
	return &record, result.Error == nil
}

func (dao *WalletDao) SaveIdempotencyKey(record *IdempotencyRecord) error {
	return dao.db.
		Table("idempotency_keys").
		Create(record).
		Error
}
