package domain

import "time"

type Transaction struct {
	ID                       uint      `gorm:"primarykey"`
	ConsumerID               uint      `gorm:"not null"`
	ConsumerCreditLimitID    uint      `gorm:"not null"`
	NomorKontrak             string    `gorm:"type:varchar(50);unique;not null"`
	TanggalKontrak           time.Time `gorm:"not null"`
	Otr                      float64   `gorm:"type:decimal(19,2);not null"`
	UangMuka                 float64   `gorm:"type:decimal(19,2);default:0"`
	AdminFee                 float64   `gorm:"type:decimal(19,2);default:0"`
	PokokPembiayaanAwal      float64   `gorm:"type:decimal(19,2);not null"`
	NilaiCicilanPerPeriode   float64   `gorm:"type:decimal(19,2);not null"`
	TenorBulan               int       `gorm:"not null"`
	TotalBunga               float64   `gorm:"type:decimal(19,2);not null"`
	TotalKewajibanPembayaran float64   `gorm:"type:decimal(19,2);not null"`
	NamaAsset                string    `gorm:"type:varchar(255)"`
	JenisAsset               string    `gorm:"type:varchar(50)"`
	StatusKontrak            string    `gorm:"type:varchar(30);not null"`
	Catatan                  string    `gorm:"type:text"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
