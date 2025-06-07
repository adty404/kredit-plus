package domain

import "time"

type ConsumerCreditLimit struct {
	ID          uint    `gorm:"primarykey"`
	ConsumerID  uint    `gorm:"not null"`
	TenorMonths int     `gorm:"not null"`
	CreditLimit float64 `gorm:"type:decimal(15,2);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
