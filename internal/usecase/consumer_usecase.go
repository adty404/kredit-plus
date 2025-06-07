package usecase

import (
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
	"time"
)

// ConsumerUsecase mendefinisikan contract untuk logika bisnis yang berhubungan dengan konsumen.
type ConsumerUsecase interface {
	CreateConsumer(input CreateConsumerInput) (*domain.Consumer, error)
	GetAllConsumers() ([]*domain.Consumer, error)
	GetConsumerByID(id uint) (*domain.Consumer, error)
	UpdateConsumer(id uint, input UpdateConsumerInput) (*domain.Consumer, error)
	DeleteConsumer(id uint) error
}

// consumerUsecase adalah implementasi dari ConsumerUsecase.
type consumerUsecase struct {
	repo domain.ConsumerRepository
}

// NewConsumerUsecase adalah factory function untuk membuat instance baru dari consumerUsecase.
// Ini menerima ConsumerRepository sebagai dependency.
func NewConsumerUsecase(repo domain.ConsumerRepository) ConsumerUsecase {
	return &consumerUsecase{repo: repo}
}

// CreateConsumer berisi logika untuk membuat konsumen baru.
func (uc *consumerUsecase) CreateConsumer(input CreateConsumerInput) (*domain.Consumer, error) {
	// Logika bisnis 1: Cek apakah NIK sudah ada.
	existingConsumer, _ := uc.repo.FindByNIK(input.Nik)
	if existingConsumer != nil {
		return nil, fmt.Errorf("consumer with NIK %s already exists", input.Nik)
	}

	// Parsing tanggal lahir dari string ke time.Time
	dob, err := time.Parse("2006-01-02", input.TanggalLahir)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for tanggal_lahir, please use yyyy-MM-dd")
	}

	jsonDob := domain.JSONDate(dob)

	// Membuat objek domain.Consumer baru dari input.
	consumer := &domain.Consumer{
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

	if err := uc.repo.Save(consumer); err != nil {
		return nil, err
	}
	return consumer, nil
}

// GetAllConsumers mengambil semua data konsumen.
func (uc *consumerUsecase) GetAllConsumers() ([]*domain.Consumer, error) {
	// Langsung memanggil repository karena tidak ada logika bisnis tambahan.
	return uc.repo.FindAll()
}

// GetConsumerByID mengambil satu konsumen berdasarkan ID.
func (uc *consumerUsecase) GetConsumerByID(id uint) (*domain.Consumer, error) {
	// Langsung memanggil repository.
	return uc.repo.FindByID(id)
}

// DeleteConsumer menghapus seorang konsumen.
func (uc *consumerUsecase) DeleteConsumer(id uint) error {
	// Cek dulu apakah konsumen ada.
	_, err := uc.repo.FindByID(id)
	if err != nil {
		return err // Mengembalikan error jika tidak ditemukan
	}

	// Panggil repository untuk menghapus.
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
