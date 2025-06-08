package seeder

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"log"
	"time"
)

// Run menjalankan semua seeder untuk mengisi data awal.
func Run(db *gorm.DB) {
	log.Println("Running database seeder...")

	if err := createBudi(db); err != nil {
		log.Fatalf("Failed to seed Budi's data: %v", err)
	}

	if err := createAnnisa(db); err != nil {
		log.Fatalf("Failed to seed Annisa's data: %v", err)
	}

	log.Println("Seeder finished successfully.")
}

// createBudi membuat data untuk konsumen Budi beserta limit kreditnya.
func createBudi(db *gorm.DB) error {
	// Tentukan tanggal lahir
	tanggalLahirBudi, _ := time.Parse("2006-01-02", "1990-05-15")

	jsonDob := domain.JSONDate(tanggalLahirBudi)

	// Data konsumen
	budi := domain.Consumer{
		Nik:                "3271011505900001", // NIK unik
		FullName:           "Budi",
		LegalName:          "Budi",
		TempatLahir:        "Bandung",
		TanggalLahir:       &jsonDob,
		Gaji:               8000000,
		OverallCreditLimit: 20000000, // Plafon kredit total Budi
	}

	if err := db.Where(domain.Consumer{Nik: budi.Nik}).FirstOrCreate(&budi).Error; err != nil {
		return err
	}
	log.Printf("Consumer '%s' created or already exists.\n", budi.FullName)

	// Data Limit Kredit untuk Budi
	var count int64
	db.Model(&domain.ConsumerCreditLimit{}).Where("consumer_id = ?", budi.ID).Count(&count)

	if count > 0 {
		log.Printf("Credit limits for '%s' already exist. Skipping.\n", budi.FullName)
		return nil
	}

	creditLimits := []domain.ConsumerCreditLimit{
		{ConsumerID: budi.ID, TenorMonths: 1, CreditLimit: 2000000},
		{ConsumerID: budi.ID, TenorMonths: 2, CreditLimit: 3500000},
		{ConsumerID: budi.ID, TenorMonths: 3, CreditLimit: 5000000},
		{ConsumerID: budi.ID, TenorMonths: 6, CreditLimit: 8000000},
	}

	// Buat semua limit kredit dalam satu batch.
	if err := db.Create(&creditLimits).Error; err != nil {
		return err
	}
	log.Printf("Successfully seeded %d credit limits for '%s'.\n", len(creditLimits), budi.FullName)

	return nil
}

// createAnnisa membuat data untuk konsumen Annisa beserta limit kreditnya.
func createAnnisa(db *gorm.DB) error {
	// Tentukan tanggal lahir
	tanggalLahirAnnisa, _ := time.Parse("2006-01-02", "1990-05-15")

	jsonDob := domain.JSONDate(tanggalLahirAnnisa)

	// Data konsumen
	annisa := domain.Consumer{
		Nik:                "3271011505900002", // NIK unik
		FullName:           "Annisa",
		LegalName:          "Annisa",
		TempatLahir:        "Bandung",
		TanggalLahir:       &jsonDob,
		Gaji:               8000000,
		OverallCreditLimit: 23200000, // Plafon kredit total Annisa
	}

	if err := db.Where(domain.Consumer{Nik: annisa.Nik}).FirstOrCreate(&annisa).Error; err != nil {
		return err
	}
	log.Printf("Consumer '%s' created or already exists.\n", annisa.FullName)

	// Data Limit Kredit untuk Annisa
	var count int64
	db.Model(&domain.ConsumerCreditLimit{}).Where("consumer_id = ?", annisa.ID).Count(&count)

	if count > 0 {
		log.Printf("Credit limits for '%s' already exist. Skipping.\n", annisa.FullName)
		return nil
	}

	creditLimits := []domain.ConsumerCreditLimit{
		{ConsumerID: annisa.ID, TenorMonths: 1, CreditLimit: 1000000},
		{ConsumerID: annisa.ID, TenorMonths: 2, CreditLimit: 1200000},
		{ConsumerID: annisa.ID, TenorMonths: 3, CreditLimit: 1500000},
		{ConsumerID: annisa.ID, TenorMonths: 6, CreditLimit: 2000000},
	}

	// Buat semua limit kredit dalam satu batch.
	if err := db.Create(&creditLimits).Error; err != nil {
		return err
	}
	log.Printf("Successfully seeded %d credit limits for '%s'.\n", len(creditLimits), annisa.FullName)

	return nil
}
