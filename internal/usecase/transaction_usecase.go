package usecase

import (
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// TransactionUsecase mendefinisikan contract untuk logika bisnis transaksi.
type TransactionUsecase interface {
	CreateTransaction(consumerID uint, input CreateTransactionInput) (*domain.Transaction, error)
	GetTransactionsByConsumerID(consumerID uint) ([]*domain.Transaction, error)
}

// transactionUsecase adalah implementasi dari TransactionUsecase.
type transactionUsecase struct {
	db              *gorm.DB // Koneksi GORM utama untuk mengelola transaksi
	transactionRepo domain.TransactionRepository
	consumerRepo    domain.ConsumerRepository
	creditLimitRepo domain.ConsumerCreditLimitRepository
}

// NewTransactionUsecase adalah factory function untuk membuat instance transactionUsecase.
func NewTransactionUsecase(
	db *gorm.DB, // Terima koneksi GORM utama
	transactionRepo domain.TransactionRepository,
	consumerRepo domain.ConsumerRepository,
	creditLimitRepo domain.ConsumerCreditLimitRepository,
) TransactionUsecase {
	return &transactionUsecase{
		db:              db,
		transactionRepo: transactionRepo,
		consumerRepo:    consumerRepo,
		creditLimitRepo: creditLimitRepo,
	}
}

func (uc *transactionUsecase) CreateTransaction(consumerID uint, input CreateTransactionInput) (
	*domain.Transaction,
	error,
) {
	var newTransaction *domain.Transaction

	// Membungkus seluruh logika dalam sebuah transaksi database.
	// Jika ada error di dalam fungsi ini, semua operasi akan di-rollback.
	err := uc.db.Transaction(
		func(tx *gorm.DB) error {
			consumerRepoTx := uc.consumerRepo.WithTx(tx)
			creditLimitRepoTx := uc.creditLimitRepo.WithTx(tx)
			transactionRepoTx := uc.transactionRepo.WithTx(tx)

			// 1. Validasi: Dapatkan data konsumen dan KUNCI barisnya untuk mencegah race condition.
			consumer, err := consumerRepoTx.FindByIDForUpdate(consumerID)
			if err != nil {
				return fmt.Errorf("consumer with id %d not found", consumerID)
			}

			// Validasi: Dapatkan limit kredit
			creditLimit, err := creditLimitRepoTx.FindByConsumerAndTenor(consumerID, input.TenorMonths)
			if err != nil {
				return fmt.Errorf("credit limit for tenor %d not found for this consumer", input.TenorMonths)
			}

			// 2. Kalkulasi Pokok Pembiayaan
			pokokPembiayaan := input.Otr - input.UangMuka + input.AdminFee

			// 3. Validasi: Cek apakah pokok pembiayaan melebihi limit produk tenor
			if pokokPembiayaan > creditLimit.CreditLimit {
				return fmt.Errorf(
					"loan amount (%.2f) exceeds tenor credit limit (%.2f)",
					pokokPembiayaan,
					creditLimit.CreditLimit,
				)
			}

			// 4. Validasi: Cek ketersediaan plafon kredit keseluruhan
			activeTransactions, err := transactionRepoTx.FindActiveByConsumerID(consumerID)
			if err != nil {
				return err
			}
			var totalPinjamanAktif float64
			for _, trans := range activeTransactions {
				totalPinjamanAktif += trans.PokokPembiayaanAwal
			}
			sisaPlafon := consumer.OverallCreditLimit - totalPinjamanAktif
			if pokokPembiayaan > sisaPlafon {
				return fmt.Errorf(
					"loan amount (%.2f) exceeds available overall credit limit (%.2f)",
					pokokPembiayaan,
					sisaPlafon,
				)
			}

			// 5. Kalkulasi detail transaksi
			totalBunga := pokokPembiayaan * 0.10
			totalKewajiban := pokokPembiayaan + totalBunga
			nilaiCicilan := totalKewajiban / float64(input.TenorMonths)

			// 6. Buat objek transaksi
			transactionToSave := &domain.Transaction{
				ConsumerID:               consumerID,
				ConsumerCreditLimitID:    creditLimit.ID,
				NomorKontrak:             fmt.Sprintf("KONTRAK/%d/%d", time.Now().Unix(), rand.Intn(1000)),
				TanggalKontrak:           time.Now(),
				Otr:                      input.Otr,
				UangMuka:                 input.UangMuka,
				AdminFee:                 input.AdminFee,
				PokokPembiayaanAwal:      pokokPembiayaan,
				NilaiCicilanPerPeriode:   nilaiCicilan,
				TenorBulan:               input.TenorMonths,
				TotalBunga:               totalBunga,
				TotalKewajibanPembayaran: totalKewajiban,
				NamaAsset:                input.NamaAsset,
				JenisAsset:               input.JenisAsset,
				StatusKontrak:            "AKTIF",
				SumberTransaksi:          input.SumberTransaksi,
			}

			// 7. Simpan transaksi
			if err = transactionRepoTx.Save(transactionToSave); err != nil {
				return err
			}
			newTransaction = transactionToSave

			// Jika tidak ada error, kembalikan nil untuk COMMIT transaksi.
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (uc *transactionUsecase) GetTransactionsByConsumerID(consumerID uint) ([]*domain.Transaction, error) {
	// Pastikan konsumen ada
	if _, err := uc.consumerRepo.FindByID(consumerID); err != nil {
		return nil, fmt.Errorf("consumer with id %d not found", consumerID)
	}
	return uc.transactionRepo.FindByConsumerID(consumerID)
}
