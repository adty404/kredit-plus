package domain

import "gorm.io/gorm"

type ConsumerCreditLimitRepository interface {
	WithTx(tx *gorm.DB) ConsumerCreditLimitRepository
	Save(creditLimit *ConsumerCreditLimit) error
	FindByConsumerAndTenor(consumerID uint, tenorMonths int) (*ConsumerCreditLimit, error)
}
