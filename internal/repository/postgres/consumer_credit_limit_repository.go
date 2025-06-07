package postgres

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
)

type consumerCreditLimitRepository struct {
	db *gorm.DB
}

func NewConsumerCreditLimitRepository(db *gorm.DB) *consumerCreditLimitRepository {
	return &consumerCreditLimitRepository{
		db: db,
	}
}

func (r *consumerCreditLimitRepository) Save(consumerCreditLimit *domain.ConsumerCreditLimit) error {
	return r.db.Create(consumerCreditLimit).Error
}

func (r *consumerCreditLimitRepository) FindByConsumerAndTenor(
	consumerID uint,
	tenorMonths int,
) (*domain.ConsumerCreditLimit, error) {
	var limit domain.ConsumerCreditLimit
	err := r.db.Where("consumer_id = ? AND tenor_months = ?", consumerID, tenorMonths).First(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}
