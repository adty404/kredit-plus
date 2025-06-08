package usecase

import (
	"errors"
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"time"
)

// ConsumerUsecase mendefinisikan contract untuk logika bisnis yang berhubungan dengan konsumen.
type ConsumerUsecase interface {
	CreateConsumer(input CreateConsumerInput) (*domain.Consumer, error)
	GetAllConsumers() ([]*domain.Consumer, error)
	GetConsumerByUserID(userID uint) (*domain.Consumer, error)
	GetConsumerByID(id uint) (*domain.Consumer, error)
	UpdateConsumer(id uint, input UpdateConsumerInput) (*domain.Consumer, error)
	DeleteConsumer(id uint) error
}

// consumerUsecase sekarang memiliki dependensi ke db dan userRepo.
type consumerUsecase struct {
	db       *gorm.DB
	repo     domain.ConsumerRepository
	userRepo domain.UserRepository
}

// NewConsumerUsecase di-update untuk menerima dependensi baru.
func NewConsumerUsecase(db *gorm.DB, repo domain.ConsumerRepository, userRepo domain.UserRepository) ConsumerUsecase {
	return &consumerUsecase{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
	}
}

// CreateConsumer sekarang membuat User dan Consumer dalam satu transaksi.
func (uc *consumerUsecase) CreateConsumer(input CreateConsumerInput) (*domain.Consumer, error) {
	var createdConsumer *domain.Consumer

	// Membungkus seluruh operasi dalam sebuah transaksi database.
	err := uc.db.Transaction(
		func(tx *gorm.DB) error {
			// Gunakan repository dengan koneksi transaksi (tx)
			userRepoTx := uc.userRepo.WithTx(tx)
			consumerRepoTx := uc.repo.WithTx(tx)

			// 1. Validasi: Cek apakah email sudah terdaftar
			_, err := userRepoTx.FindByEmail(input.Email)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("email '%s' already registered", input.Email)
			}

			// 2. Validasi: Cek apakah NIK sudah terdaftar
			_, err = consumerRepoTx.FindByNIK(input.Nik)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("consumer with NIK %s already exists", input.Nik)
			}

			// 3. Buat User baru
			newUser := &domain.User{
				FullName: input.FullName,
				Email:    input.Email,
				Role:     "consumer",
			}
			if err := newUser.HashPassword(input.Password); err != nil {
				return err
			}
			if err := userRepoTx.Save(newUser); err != nil {
				return err
			}

			// 4. Buat Consumer baru dan tautkan UserID
			dob, err := time.Parse("2006-01-02", input.TanggalLahir)
			if err != nil {
				return fmt.Errorf("invalid date format for tanggal_lahir, please use yyyy-MM-dd")
			}
			jsonDob := domain.JSONDate(dob)

			consumer := &domain.Consumer{
				UserID:             newUser.ID, // Link ke user yang baru dibuat
				Nik:                input.Nik,
				FullName:           input.FullName,
				LegalName:          input.LegalName,
				TempatLahir:        input.TempatLahir,
				TanggalLahir:       &jsonDob,
				Gaji:               input.Gaji,
				OverallCreditLimit: input.OverallCreditLimit,
				FotoKtp:            input.FotoKtpPath,
				FotoSelfie:         input.FotoSelfiePath,
			}

			if err := consumerRepoTx.Save(consumer); err != nil {
				return err
			}

			createdConsumer = consumer
			return nil // Commit transaksi jika tidak ada error
		},
	)

	if err != nil {
		return nil, err
	}

	return createdConsumer, nil
}

// GetAllConsumers mengambil semua data konsumen.
func (uc *consumerUsecase) GetAllConsumers() ([]*domain.Consumer, error) {
	return uc.repo.FindAll()
}

// GetConsumerByUserID mengambil konsumen berdasarkan ID pengguna.
func (uc *consumerUsecase) GetConsumerByUserID(userID uint) (*domain.Consumer, error) {
	consumer, err := uc.repo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("consumer not found for user ID %d: %w", userID, err)
	}
	return consumer, nil
}

// GetConsumerByID mengambil satu konsumen berdasarkan ID.
func (uc *consumerUsecase) GetConsumerByID(id uint) (*domain.Consumer, error) {
	return uc.repo.FindByID(id)
}

// DeleteConsumer menghapus seorang konsumen.
func (uc *consumerUsecase) DeleteConsumer(id uint) error {
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return err // Mengembalikan error jika tidak ditemukan
	}

	return uc.repo.Delete(id)
}

// UpdateConsumer memperbarui data konsumen yang ada.
func (uc *consumerUsecase) UpdateConsumer(id uint, input UpdateConsumerInput) (*domain.Consumer, error) {
	// Pertama, pastikan konsumennya ada.
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Buat map untuk menampung field yang akan diupdate.
	updates := make(map[string]interface{})

	if input.FullName != nil {
		updates["full_name"] = *input.FullName
	}
	if input.LegalName != nil {
		updates["legal_name"] = *input.LegalName
	}
	if input.TempatLahir != nil {
		updates["tempat_lahir"] = *input.TempatLahir
	}
	if input.Gaji != nil {
		updates["gaji"] = *input.Gaji
	}
	if input.TanggalLahir != nil {
		dob, err := time.Parse("2006-01-02", *input.TanggalLahir)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for tanggal_lahir, please use yyyy-MM-dd")
		}
		updates["tanggal_lahir"] = dob
	}
	if input.OverallCreditLimit != nil {
		updates["overall_credit_limit"] = *input.OverallCreditLimit
	}
	if input.FotoKtp != nil {
		updates["foto_ktp"] = *input.FotoKtp
	}
	if input.FotoSelfie != nil {
		updates["foto_selfie"] = *input.FotoSelfie
	}

	// Hanya jalankan update jika ada data yang perlu diubah.
	if len(updates) > 0 {
		if err := uc.repo.Update(id, updates); err != nil {
			return nil, err
		}
	}

	// Setelah update berhasil, ambil kembali data terbaru untuk dikembalikan.
	return uc.repo.FindByID(id)
}
