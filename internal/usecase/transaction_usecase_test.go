package usecase

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// setupMocksAndDb membuat instance mock untuk semua repository dan db GORM palsu untuk testing.
func setupMocksAndDb(t *testing.T) (
	*gorm.DB,
	sqlmock.Sqlmock,
	*MockConsumerRepository,
	*MockCreditLimitRepository,
	*MockTransactionRepository,
) {
	sqlDB, mock, err := sqlmock.New()
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
	mockLimitRepo := new(MockCreditLimitRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	return gormDB, mock, mockConsumerRepo, mockLimitRepo, mockTransactionRepo
}

func TestCreateTransaction_Success(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockLimitRepo, mockTransactionRepo := setupMocksAndDb(t)
	usecase := NewTransactionUsecase(gormDB, mockTransactionRepo, mockConsumerRepo, mockLimitRepo)

	consumerID := uint(1)
	input := CreateTransactionInput{
		TenorMonths: 6,
		Otr:         5000000,
		AdminFee:    100000,
		UangMuka:    500000,
	}

	consumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 10000000}
	creditLimit := &domain.ConsumerCreditLimit{ID: 10, ConsumerID: consumerID, CreditLimit: 5000000}
	activeTransactions := []*domain.Transaction{}

	// Tentukan ekspektasi untuk transaksi SQL
	mockSQL.ExpectBegin()

	// Tentukan ekspektasi untuk metode repository yang sebenarnya
	mockConsumerRepo.On("FindByIDForUpdate", consumerID).Return(consumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(creditLimit, nil).Once()
	mockTransactionRepo.On("FindActiveByConsumerID", consumerID).Return(activeTransactions, nil).Once()
	mockTransactionRepo.On("Save", mock.AnythingOfType("*domain.Transaction")).Return(nil).Once()

	// Harapkan Commit setelah semua operasi berhasil
	mockSQL.ExpectCommit()

	// Act
	transaction, err := usecase.CreateTransaction(consumerID, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(4600000), transaction.PokokPembiayaanAwal)

	// Verifikasi semua ekspektasi (termasuk SQL) terpenuhi
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestCreateTransaction_ExceedsOverallLimit(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockLimitRepo, mockTransactionRepo := setupMocksAndDb(t)
	usecase := NewTransactionUsecase(gormDB, mockTransactionRepo, mockConsumerRepo, mockLimitRepo)

	consumerID := uint(1)
	input := CreateTransactionInput{TenorMonths: 6, Otr: 5000000}

	consumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 10000000}
	creditLimit := &domain.ConsumerCreditLimit{ID: 10, ConsumerID: consumerID, CreditLimit: 8000000}
	activeTransactions := []*domain.Transaction{{PokokPembiayaanAwal: 6000000}}

	// Tentukan ekspektasi SQL (gagal, jadi akan di-rollback)
	mockSQL.ExpectBegin()
	mockSQL.ExpectRollback()

	// Tentukan ekspektasi mock repository
	mockConsumerRepo.On("FindByIDForUpdate", consumerID).Return(consumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(creditLimit, nil).Once()
	mockTransactionRepo.On("FindActiveByConsumerID", consumerID).Return(activeTransactions, nil).Once()

	// Act
	transaction, err := usecase.CreateTransaction(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "exceeds available overall credit limit")

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestCreateTransaction_ExceedsTenorLimit(t *testing.T) {
	// Arrange
	gormDB, mockSQL, mockConsumerRepo, mockLimitRepo, mockTransactionRepo := setupMocksAndDb(t)
	usecase := NewTransactionUsecase(gormDB, mockTransactionRepo, mockConsumerRepo, mockLimitRepo)

	consumerID := uint(1)
	input := CreateTransactionInput{TenorMonths: 3, Otr: 6000000}

	consumer := &domain.Consumer{ID: consumerID, OverallCreditLimit: 10000000}
	creditLimit := &domain.ConsumerCreditLimit{ID: 10, ConsumerID: consumerID, CreditLimit: 5000000}

	// Tentukan ekspektasi SQL (gagal, jadi akan di-rollback)
	mockSQL.ExpectBegin()
	mockSQL.ExpectRollback()

	// Tentukan ekspektasi mock repository
	mockConsumerRepo.On("FindByIDForUpdate", consumerID).Return(consumer, nil).Once()
	mockLimitRepo.On("FindByConsumerAndTenor", consumerID, input.TenorMonths).Return(creditLimit, nil).Once()

	// Act
	transaction, err := usecase.CreateTransaction(consumerID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, transaction)
	assert.Contains(t, err.Error(), "exceeds tenor credit limit")

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	mockConsumerRepo.AssertExpectations(t)
	mockLimitRepo.AssertExpectations(t)
}
