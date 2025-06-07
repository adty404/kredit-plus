package usecase

import (
	"errors"
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestCreateConsumerCreditLimit_Success(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(1)
	input := CreateConsumerCreditLimitInput{
		TenorMonths: 6,
		CreditLimit: 10000000,
	}

	// Mock data konsumen yang ada
	existingConsumer := &domain.Consumer{
		ID:                 consumerID,
		OverallCreditLimit: 15000000,
	}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(existingConsumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(nil, gorm.ErrRecordNotFound).Once()
	mockLimitRepo.On("Save", mock.AnythingOfType("*domain.ConsumerCreditLimit")).Return(nil).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, limit)
	assert.Equal(t, input.CreditLimit, limit.CreditLimit)
	assert.Equal(t, input.TenorMonths, limit.TenorMonths)
	assert.Equal(t, consumerID, limit.ConsumerID)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestCreateConsumerCreditLimit_ConsumerNotFound(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(99) // ID yang tidak ada
	input := CreateConsumerCreditLimitInput{TenorMonths: 6, CreditLimit: 10000000}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, limit)
	assert.Equal(t, fmt.Errorf("consumer with id %d not found", consumerID), err)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertNotCalled(t, "FindByConsumerAndTenor")
}

func TestCreateConsumerCreditLimit_LimitAlreadyExists(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(1)
	input := CreateConsumerCreditLimitInput{TenorMonths: 6, CreditLimit: 10000000}

	existingConsumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 15000000}
	existingLimit := &domain.ConsumerCreditLimit{ID: 10, ConsumerID: consumerID, TenorMonths: input.TenorMonths}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(existingConsumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(existingLimit, nil).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, limit)
	assert.Equal(
		t,
		fmt.Errorf("credit limit for tenor %d months already exists for this consumer", input.TenorMonths),
		err,
	)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestCreateConsumerCreditLimit_InvalidTenor(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(1)
	input := CreateConsumerCreditLimitInput{TenorMonths: 5, CreditLimit: 10000000} // Tenor 5 tidak valid

	existingConsumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 15000000}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(existingConsumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, limit)
	assert.Equal(t, fmt.Errorf("invalid tenor: %d. allowed tenors are 1, 2, 3, 6", input.TenorMonths), err)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestCreateConsumerCreditLimit_ExceedsOverallLimit(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(1)
	input := CreateConsumerCreditLimitInput{TenorMonths: 6, CreditLimit: 20000000} // Melebihi overall limit

	existingConsumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 15000000}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(existingConsumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, limit)
	assert.Equal(
		t,
		fmt.Errorf(
			"credit limit (%.2f) cannot exceed consumer's overall credit limit (%.2f)",
			input.CreditLimit,
			existingConsumer.OverallCreditLimit,
		),
		err,
	)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}

func TestCreateConsumerCreditLimit_SaveError(t *testing.T) {
	// Arrange
	mockConsumerRepo := new(MockConsumerRepository)
	mockLimitRepo := new(MockCreditLimitRepository)
	usecase := NewConsumerCreditLimitUsecase(mockLimitRepo, mockConsumerRepo)

	consumerID := uint(1)
	input := CreateConsumerCreditLimitInput{TenorMonths: 6, CreditLimit: 10000000}
	dbError := errors.New("database save error")

	existingConsumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 15000000}

	// Tentukan ekspektasi
	mockConsumerRepo.On("FindByID", consumerID).Return(existingConsumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(nil, gorm.ErrRecordNotFound).Once()
	mockLimitRepo.On("Save", mock.AnythingOfType("*domain.ConsumerCreditLimit")).Return(dbError).Once()

	// Act
	limit, err := usecase.CreateConsumerCreditLimit(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, limit)
	assert.Equal(t, dbError, err)
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}
