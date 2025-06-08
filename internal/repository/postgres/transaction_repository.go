package postgres

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) WithTx(tx *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db: tx}
}

func (r *transactionRepository) Save(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	if err := r.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByConsumerID(consumerID uint) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	if err := r.db.Where(
		"consumer_id = ?",
		consumerID,
	).Order("tanggal_kontrak desc").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) FindActiveByConsumerID(consumerID uint) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.Where("consumer_id = ? AND status_kontrak = ?", consumerID, "AKTIF").Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	return r.db.Save(transaction).Error
}
