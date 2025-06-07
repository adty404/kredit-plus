package domain

import "time"

type Consumer struct {
	ID                 uint      `gorm:"primarykey"`
	Nik                string    `gorm:"type:varchar(16);unique;not null"`
	FullName           string    `gorm:"type:varchar(255);not null"`
	LegalName          string    `gorm:"type:varchar(255)"`
	TempatLahir        string    `gorm:"type:varchar(100)"`
	TanggalLahir       *JSONDate `gorm:"type:date"`
	Gaji               float64   `gorm:"type:decimal(15,2)"`
	OverallCreditLimit float64   `gorm:"type:decimal(19,2);not null;default:0"`
	FotoKtp            string    `gorm:"type:varchar(255)"`
	FotoSelfie         string    `gorm:"type:varchar(255)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relasi
	CreditLimits []ConsumerCreditLimit `gorm:"foreignKey:ConsumerID"`
	Transactions []Transaction         `gorm:"foreignKey:ConsumerID"`
}
