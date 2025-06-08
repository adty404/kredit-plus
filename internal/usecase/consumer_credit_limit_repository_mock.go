package usecase

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockCreditLimitRepository adalah implementasi mock dari domain.ConsumerCreditLimitRepository.
type MockCreditLimitRepository struct {
	mock.Mock
}

func (m *MockCreditLimitRepository) WithTx(tx *gorm.DB) domain.ConsumerCreditLimitRepository {
	return m
}

func (m *MockCreditLimitRepository) Save(limit *domain.ConsumerCreditLimit) error {
	args := m.Called(limit)
	return args.Error(0)
}

func (m *MockCreditLimitRepository) Update(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockCreditLimitRepository) FindByID(id uint) (*domain.ConsumerCreditLimit, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ConsumerCreditLimit), args.Error(1)
}

func (m *MockCreditLimitRepository) FindByConsumerID(consumerID uint) ([]*domain.ConsumerCreditLimit, error) {
	args := m.Called(consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ConsumerCreditLimit), args.Error(1)
}

func (m *MockCreditLimitRepository) FindByConsumerAndTenor(
	consumerID uint,
	tenorMonths int,
) (*domain.ConsumerCreditLimit, error) {
	args := m.Called(consumerID, tenorMonths)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ConsumerCreditLimit), args.Error(1)
}

func (m *MockCreditLimitRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
