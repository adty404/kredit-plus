package domain

type ConsumerCreditLimitRepository interface {
	Save(creditLimit *ConsumerCreditLimit) error
	FindByConsumerAndTenor(consumerID uint, tenorMonths int) (*ConsumerCreditLimit, error)
}
