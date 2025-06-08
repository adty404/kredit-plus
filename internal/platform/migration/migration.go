package migration

import (
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
	"gorm.io/gorm"
	"log"
)

func Run(db *gorm.DB) error {
	log.Println("Starting GORM auto-migration...")

	err := db.AutoMigrate(
		&domain.Consumer{},
		&domain.ConsumerCreditLimit{},
		&domain.Transaction{},
		&domain.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	log.Println("Database migration completed successfully!")
	return nil
}
