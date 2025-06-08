package usecase

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupMocksForConsumerTest adalah helper untuk membuat semua mock yang dibutuhkan.
func setupMocksForConsumerTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *MockConsumerRepository, *MockUserRepository) {
	sqlDB, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(
		postgres.New(
			postgres.Config{
				Conn: sqlDB,
			},
		), &gorm.Config{},
	)
	assert.NoError(t, err)

	mockConsumerRepo := new(MockConsumerRepository)
	mockUserRepo := new(MockUserRepository)

	return gormDB, mockSQL, mockConsumerRepo, mockUserRepo
}

// --- Test untuk CreateConsumer ---

func TestConsumerUsecase_CreateConsumer_Success(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)

	input := CreateConsumerInput{
		Nik:          "1234567890123456",
		FullName:     "Test Consumer",
		Email:        "consumer@example.com",
		Password:     "password123",
		TanggalLahir: "2000-01-01",
	}

	mockSQL.ExpectBegin()
	mockUserRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()
	mockConsumerRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()
	mockUserRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()
	mockConsumerRepo.On("Save", mock.AnythingOfType("*domain.Consumer")).Return(nil).Once()
	mockSQL.ExpectCommit()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, input.Nik, consumer.Nik)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockConsumerRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestConsumerUsecase_CreateConsumer_NikExists(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	input := CreateConsumerInput{Nik: "123", Email: "new@example.com"}

	mockSQL.ExpectBegin()
	mockUserRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()
	mockConsumerRepo.On("FindByNIK", input.Nik).Return(&domain.Consumer{}, nil).Once()
	mockSQL.ExpectRollback()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Contains(t, err.Error(), "consumer with NIK 123 already exists")
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockUserRepo.AssertExpectations(t)
	mockConsumerRepo.AssertExpectations(t)
}

func TestConsumerUsecase_CreateConsumer_SaveError(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	input := CreateConsumerInput{Nik: "123", Email: "new@example.com", TanggalLahir: "2000-01-01"}
	dbError := errors.New("database save error")

	mockSQL.ExpectBegin()
	mockUserRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()
	mockConsumerRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()
	mockUserRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()
	mockConsumerRepo.On(
		"Save",
		mock.AnythingOfType("*domain.Consumer"),
	).Return(dbError).Once() // Simulasikan error di sini
	mockSQL.ExpectRollback()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, dbError, err)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockConsumerRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestConsumerUsecase_CreateConsumer_InvalidDateFormat(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	input := CreateConsumerInput{Nik: "123", Email: "new@example.com", TanggalLahir: "01-01-2000"} // Format salah

	mockSQL.ExpectBegin()
	mockUserRepo.On("FindByEmail", input.Email).Return(nil, gorm.ErrRecordNotFound).Once()
	mockConsumerRepo.On("FindByNIK", input.Nik).Return(nil, gorm.ErrRecordNotFound).Once()
	mockUserRepo.On("Save", mock.AnythingOfType("*domain.User")).Return(nil).Once()
	mockSQL.ExpectRollback()

	// Act
	consumer, err := usecase.CreateConsumer(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Contains(t, err.Error(), "invalid date format")
	assert.NoError(t, mockSQL.ExpectationsWereMet())
}

// --- Test untuk GetConsumerByID ---

func TestConsumerUsecase_GetConsumerByID_Success(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	expectedConsumer := &domain.Consumer{ID: 1, FullName: "Test User"}

	mockConsumerRepo.On("FindByID", uint(1)).Return(expectedConsumer, nil).Once()

	// Act
	consumer, err := usecase.GetConsumerByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, expectedConsumer.ID, consumer.ID)
	mockConsumerRepo.AssertExpectations(t)
}

func TestConsumerUsecase_GetConsumerByID_NotFound(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)

	mockConsumerRepo.On("FindByID", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	consumer, err := usecase.GetConsumerByID(1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockConsumerRepo.AssertExpectations(t)
}

// --- Test untuk GetAllConsumers ---

func TestConsumerUsecase_GetAllConsumers_Success(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	expectedConsumers := []*domain.Consumer{
		{ID: 1, FullName: "User Satu"},
		{ID: 2, FullName: "User Dua"},
	}

	mockConsumerRepo.On("FindAll").Return(expectedConsumers, nil).Once()

	// Act
	consumers, err := usecase.GetAllConsumers()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, consumers, 2)
	mockConsumerRepo.AssertExpectations(t)
}

func TestConsumerUsecase_GetAllConsumers_Empty(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	expectedConsumers := []*domain.Consumer{}

	mockConsumerRepo.On("FindAll").Return(expectedConsumers, nil).Once()

	// Act
	consumers, err := usecase.GetAllConsumers()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, consumers, 0)
	mockConsumerRepo.AssertExpectations(t)
}

// --- Test untuk UpdateConsumer ---

func TestConsumerUsecase_UpdateConsumer_Success(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	idToUpdate := uint(1)
	newName := "Updated Name"
	input := UpdateConsumerInput{FullName: &newName}
	initialConsumer := &domain.Consumer{ID: idToUpdate, FullName: "Initial Name"}
	updatedConsumer := &domain.Consumer{ID: idToUpdate, FullName: *input.FullName}

	mockConsumerRepo.On("FindByID", idToUpdate).Return(initialConsumer, nil).Once()
	mockConsumerRepo.On("Update", idToUpdate, mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()
	mockConsumerRepo.On("FindByID", idToUpdate).Return(updatedConsumer, nil).Once()

	// Act
	consumer, err := usecase.UpdateConsumer(idToUpdate, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", consumer.FullName)
	mockConsumerRepo.AssertExpectations(t)
}

func TestConsumerUsecase_UpdateConsumer_NotFound(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	idToUpdate := uint(99)
	newName := "Updated Name"
	input := UpdateConsumerInput{FullName: &newName}

	mockConsumerRepo.On("FindByID", idToUpdate).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	consumer, err := usecase.UpdateConsumer(idToUpdate, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockConsumerRepo.AssertExpectations(t)
}

// --- Test untuk DeleteConsumer ---

func TestConsumerUsecase_DeleteConsumer_Success(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	idToDelete := uint(1)

	mockConsumerRepo.On("FindByID", idToDelete).Return(&domain.Consumer{ID: idToDelete}, nil).Once()
	mockConsumerRepo.On("Delete", idToDelete).Return(nil).Once()

	// Act
	err := usecase.DeleteConsumer(idToDelete)

	// Assert
	assert.NoError(t, err)
	mockConsumerRepo.AssertExpectations(t)
}

func TestConsumerUsecase_DeleteConsumer_NotFound(t *testing.T) {
	// Arrange
	gormDB, _, mockConsumerRepo, mockUserRepo := setupMocksForConsumerTest(t)
	usecase := NewConsumerUsecase(gormDB, mockConsumerRepo, mockUserRepo)
	idToDelete := uint(99)

	mockConsumerRepo.On("FindByID", idToDelete).Return(nil, gorm.ErrRecordNotFound).Once()

	// Act
	err := usecase.DeleteConsumer(idToDelete)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	mockConsumerRepo.AssertExpectations(t)
}
