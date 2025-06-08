package usecase

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTransactionRepository adalah implementasi mock dari domain.TransactionRepository.
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) WithTx(tx *gorm.DB) domain.TransactionRepository {
	return m
}

func (m *MockTransactionRepository) Save(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(id uint) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByConsumerID(consumerID uint) ([]*domain.Transaction, error) {
	args := m.Called(consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindActiveByConsumerID(consumerID uint) ([]*domain.Transaction, error) {
	args := m.Called(consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}
