package domain

import "gorm.io/gorm"

type TransactionRepository interface {
	WithTx(tx *gorm.DB) TransactionRepository
	Save(transaction *Transaction) error
	FindByID(id uint) (*Transaction, error)
	FindByConsumerID(consumerID uint) ([]*Transaction, error)
	// FindActiveByConsumerID hanya mengambil transaksi dengan status 'AKTIF' untuk perhitungan limit.
	FindActiveByConsumerID(consumerID uint) ([]*Transaction, error)
	Update(transaction *Transaction) error
}
