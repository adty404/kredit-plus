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

// --- Test untuk CreateConsumer ---

func TestCreateConsumer_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	input := CreateConsumerInput{
		Nik:          "1234567890123456",
		FullName:     "Test User",
		LegalName:    "Test User Legal",
		TempatLahir:  "Test City",
		TanggalLahir: "2000-01-01",
		Gaji:         10000000,
	}

	mockRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("*domain.Consumer")).Return(nil).Once()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, input.Nik, consumer.Nik)
	mockRepo.AssertExpectations(t)
}

func TestCreateConsumer_NikExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	input := CreateConsumerInput{Nik: "1234567890123456"}
	existingConsumer := &domain.Consumer{ID: 1, Nik: input.Nik}

	mockRepo.On("FindByNIK", input.Nik).Return(existingConsumer, nil).Once()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, fmt.Errorf("consumer with NIK %s already exists", input.Nik), err)
	mockRepo.AssertExpectations(t)
}

func TestCreateConsumer_SaveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	input := CreateConsumerInput{Nik: "1234567890123456", TanggalLahir: "2000-01-01"}
	dbError := errors.New("database save error")

	mockRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()
	mockRepo.On("Save", mock.AnythingOfType("*domain.Consumer")).Return(dbError).Once()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, dbError, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateConsumer_InvalidDateFormat(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	input := CreateConsumerInput{Nik: "1234567890123456", TanggalLahir: "01-01-2000"} // Format salah

	// Tambahkan ekspektasi untuk FindByNIK karena dipanggil sebelum validasi tanggal.
	// Asumsikan NIK tidak ditemukan untuk melanjutkan ke validasi tanggal.
	mockRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Contains(t, err.Error(), "invalid date format")
	mockRepo.AssertExpectations(t)
}

// --- Test untuk GetConsumerByID ---

func TestGetConsumerByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	expectedConsumer := &domain.Consumer{ID: 1, FullName: "Test User"}

	mockRepo.On("FindByID", uint(1)).Return(expectedConsumer, nil).Once()

	// Act
	consumer, err := usecase.GetConsumerByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, expectedConsumer.ID, consumer.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetConsumerByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	consumer, err := usecase.GetConsumerByID(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockRepo.AssertExpectations(t)
}

// --- Test untuk GetAllConsumers ---

func TestGetAllConsumers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	expectedConsumers := []*domain.Consumer{
		{ID: 1, FullName: "User Satu"},
		{ID: 2, FullName: "User Dua"},
	}

	mockRepo.On("FindAll").Return(expectedConsumers, nil).Once()

	// Act
	consumers, err := usecase.GetAllConsumers()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumers)
	assert.Len(t, consumers, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetAllConsumers_Empty(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	// Mengembalikan slice kosong
	expectedConsumers := []*domain.Consumer{}

	mockRepo.On("FindAll").Return(expectedConsumers, nil).Once()

	// Act
	consumers, err := usecase.GetAllConsumers()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumers) // Slice tidak nil, tapi kosong
	assert.Len(t, consumers, 0)
	mockRepo.AssertExpectations(t)
}

// --- Test untuk UpdateConsumer ---

func TestUpdateConsumer_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)

	idToUpdate := uint(1)
	newName := "Updated Name"
	input := UpdateConsumerInput{FullName: &newName}

	initialConsumer := &domain.Consumer{ID: idToUpdate, FullName: "Initial Name"}
	updatedConsumer := &domain.Consumer{ID: idToUpdate, FullName: *input.FullName}

	mockRepo.On("FindByID", idToUpdate).Return(initialConsumer, nil).Once()
	mockRepo.On("Update", idToUpdate, mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()
	mockRepo.On("FindByID", idToUpdate).Return(updatedConsumer, nil).Once()

	// Act
	consumer, err := usecase.UpdateConsumer(idToUpdate, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, *input.FullName, consumer.FullName)
	mockRepo.AssertExpectations(t)
}

func TestUpdateConsumer_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	idToUpdate := uint(99) // ID yang tidak ada
	newName := "Updated Name"
	input := UpdateConsumerInput{FullName: &newName}

	mockRepo.On("FindByID", idToUpdate).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	consumer, err := usecase.UpdateConsumer(idToUpdate, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockRepo.AssertExpectations(t)
}

// --- Test untuk DeleteConsumer ---

func TestDeleteConsumer_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	idToDelete := uint(1)

	mockRepo.On("FindByID", idToDelete).Return(&domain.Consumer{ID: idToDelete}, nil).Once()
	mockRepo.On("Delete", idToDelete).Return(nil).Once()

	// Act
	err := usecase.DeleteConsumer(idToDelete)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteConsumer_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockConsumerRepository)
	usecase := NewConsumerUsecase(mockRepo)
	idToDelete := uint(99)

	mockRepo.On("FindByID", idToDelete).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	err := usecase.DeleteConsumer(idToDelete)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockRepo.AssertExpectations(t)
}
