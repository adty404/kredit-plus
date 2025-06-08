package usecase

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockConsumerRepository struct {
	mock.Mock
}

func (m *MockConsumerRepository) WithTx(tx *gorm.DB) domain.ConsumerRepository {
	return m
}

func (m *MockConsumerRepository) FindByIDForUpdate(id uint) (*domain.Consumer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Consumer), args.Error(1)
}

// FindByUserID adalah metode baru yang ditambahkan untuk memenuhi interface.
func (m *MockConsumerRepository) FindByUserID(userID uint) (*domain.Consumer, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Consumer), args.Error(1)
}

func (m *MockConsumerRepository) Save(consumer *domain.Consumer) error {
	args := m.Called(consumer)
	return args.Error(0)
}

func (m *MockConsumerRepository) Update(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockConsumerRepository) FindByID(id uint) (*domain.Consumer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Consumer), args.Error(1)
}

func (m *MockConsumerRepository) FindByNIK(nik string) (*domain.Consumer, error) {
	args := m.Called(nik)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Consumer), args.Error(1)
}

func (m *MockConsumerRepository) FindAll() ([]*domain.Consumer, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Consumer), args.Error(1)
}

func (m *MockConsumerRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
