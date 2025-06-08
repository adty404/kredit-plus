package seeder

import (
	"errors"
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"log"
	"time"
)

// Run menjalankan semua seeder untuk mengisi data awal.
func Run(db *gorm.DB) {
	log.Println("Running database seeder...")

	// Jalankan seeder untuk admin user
	if err := createAdminUser(db); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}

	// Jalankan seeder untuk konsumen Budi (user dan data consumer)
	if err := createBudi(db); err != nil {
		log.Fatalf("Failed to seed Budi's data: %v", err)
	}

	// Jalankan seeder untuk konsumen Annisa (user dan data consumer)
	if err := createAnnisa(db); err != nil {
		log.Fatalf("Failed to seed Annisa's data: %v", err)
	}

	log.Println("Seeder finished successfully.")
}

// createAdminUser membuat pengguna dengan peran 'admin' jika belum ada.
func createAdminUser(db *gorm.DB) error {
	adminEmail := "admin@kreditplus.com"
	var existingUser domain.User
	err := db.Where("email = ?", adminEmail).First(&existingUser).Error
	if err == nil {
		log.Printf("Admin user with email '%s' already exists. Skipping.\n", adminEmail)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	adminUser := &domain.User{
		FullName: "Admin Kredit Plus",
		Email:    adminEmail,
		Role:     "admin",
	}
	if err := adminUser.HashPassword("password123"); err != nil {
		return err
	}
	if err := db.Create(adminUser).Error; err != nil {
		return err
	}
	log.Printf("Successfully seeded admin user with email '%s'\n", adminEmail)
	return nil
}

// createUserIfNotExists adalah helper untuk membuat user dan mengembalikan ID-nya.
func createUserIfNotExists(db *gorm.DB, fullName, email, password, role string) (uint, error) {
	var existingUser domain.User
	err := db.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		log.Printf("User with email '%s' already exists. Using existing ID.\n", email)
		return existingUser.ID, nil // Kembalikan ID user yang sudah ada
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	newUser := &domain.User{
		FullName: fullName,
		Email:    email,
		Role:     role,
	}
	if err := newUser.HashPassword(password); err != nil {
		return 0, err
	}
	if err := db.Create(newUser).Error; err != nil {
		return 0, err
	}
	log.Printf("Successfully seeded user '%s'.\n", fullName)
	return newUser.ID, nil
}

// createBudi membuat data user dan consumer untuk Budi.
func createBudi(db *gorm.DB) error {
	// 1. Buat User untuk Budi
	budiUserID, err := createUserIfNotExists(db, "Budi Santoso", "budi@example.com", "passwordbudi", "consumer")
	if err != nil {
		return err
	}

	// 2. Buat Consumer Budi dan tautkan dengan UserID
	var consumerBudi domain.Consumer
	err = db.Where("nik = ?", "3271011505900001").First(&consumerBudi).Error
	if err == nil {
		log.Printf("Consumer 'Budi' already exists. Skipping consumer creation.\n")
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		tanggalLahirBudi, _ := time.Parse("2006-01-02", "1990-05-15")
		jsonDob := domain.JSONDate(tanggalLahirBudi)
		consumerBudi = domain.Consumer{
			UserID:             budiUserID,
			Nik:                "3271011505900001",
			FullName:           "Budi Santoso",
			LegalName:          "Budi Santoso",
			TempatLahir:        "Bandung",
			TanggalLahir:       &jsonDob,
			Gaji:               8000000,
			OverallCreditLimit: 20000000,
		}
		if err := db.Create(&consumerBudi).Error; err != nil {
			return err
		}
		log.Printf("Consumer '%s' created.\n", consumerBudi.FullName)
	} else {
		return err
	}

	// 3. Buat Limit Kredit untuk Budi
	var count int64
	db.Model(&domain.ConsumerCreditLimit{}).Where("consumer_id = ?", consumerBudi.ID).Count(&count)
	if count > 0 {
		log.Printf("Credit limits for '%s' already exist. Skipping.\n", consumerBudi.FullName)
		return nil
	}
	creditLimits := []domain.ConsumerCreditLimit{
		{ConsumerID: consumerBudi.ID, TenorMonths: 1, CreditLimit: 2000000},
		{ConsumerID: consumerBudi.ID, TenorMonths: 2, CreditLimit: 3500000},
		{ConsumerID: consumerBudi.ID, TenorMonths: 3, CreditLimit: 5000000},
		{ConsumerID: consumerBudi.ID, TenorMonths: 6, CreditLimit: 8000000},
	}
	if err := db.Create(&creditLimits).Error; err != nil {
		return err
	}
	log.Printf("Successfully seeded %d credit limits for '%s'.\n", len(creditLimits), consumerBudi.FullName)
	return nil
}

// createAnnisa membuat data user dan consumer untuk Annisa.
func createAnnisa(db *gorm.DB) error {
	// 1. Buat User untuk Annisa
	annisaUserID, err := createUserIfNotExists(
		db,
		"Annisa Fitriani",
		"annisa@example.com",
		"passwordannisa",
		"consumer",
	)
	if err != nil {
		return err
	}

	// 2. Buat Consumer Annisa dan tautkan dengan UserID
	var consumerAnnisa domain.Consumer
	err = db.Where("nik = ?", "3271011505900002").First(&consumerAnnisa).Error
	if err == nil {
		log.Printf("Consumer 'Annisa' already exists. Skipping consumer creation.\n")
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		tanggalLahirAnnisa, _ := time.Parse("2006-01-02", "1992-08-20")
		jsonDob := domain.JSONDate(tanggalLahirAnnisa)
		consumerAnnisa = domain.Consumer{
			UserID:             annisaUserID,
			Nik:                "3271011505900002",
			FullName:           "Annisa Fitriani",
			LegalName:          "Annisa Fitriani",
			TempatLahir:        "Jakarta",
			TanggalLahir:       &jsonDob,
			Gaji:               12000000,
			OverallCreditLimit: 25000000,
		}
		if err := db.Create(&consumerAnnisa).Error; err != nil {
			return err
		}
		log.Printf("Consumer '%s' created.\n", consumerAnnisa.FullName)
	} else {
		return err
	}

	// 3. Buat Limit Kredit untuk Annisa
	var count int64
	db.Model(&domain.ConsumerCreditLimit{}).Where("consumer_id = ?", consumerAnnisa.ID).Count(&count)
	if count > 0 {
		log.Printf("Credit limits for '%s' already exist. Skipping.\n", consumerAnnisa.FullName)
		return nil
	}
	creditLimits := []domain.ConsumerCreditLimit{
		{ConsumerID: consumerAnnisa.ID, TenorMonths: 1, CreditLimit: 1000000},
		{ConsumerID: consumerAnnisa.ID, TenorMonths: 2, CreditLimit: 1200000},
		{ConsumerID: consumerAnnisa.ID, TenorMonths: 3, CreditLimit: 1500000},
		{ConsumerID: consumerAnnisa.ID, TenorMonths: 6, CreditLimit: 2000000},
	}
	if err := db.Create(&creditLimits).Error; err != nil {
		return err
	}
	log.Printf("Successfully seeded %d credit limits for '%s'.\n", len(creditLimits), consumerAnnisa.FullName)
	return nil
}
