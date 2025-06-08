package usecase

import "mime/multipart"

type CreateConsumerFormInput struct {
	Nik                string                `form:"nik" binding:"required,len=16"`
	FullName           string                `form:"full_name" binding:"required,min=2"`
	LegalName          string                `form:"legal_name" binding:"required,min=2"`
	Email              string                `form:"email" binding:"required,email"`
	Password           string                `form:"password" binding:"required,min=8"`
	TempatLahir        string                `form:"tempat_lahir" binding:"required"`
	TanggalLahir       string                `form:"tanggal_lahir" binding:"required,datetime=2006-01-02"`
	Gaji               string                `form:"gaji" binding:"required,numeric,gt=0"`
	OverallCreditLimit string                `form:"overall_credit_limit" binding:"required,numeric,gte=0"`
	FotoKtp            *multipart.FileHeader `form:"foto_ktp" binding:"omitempty"`
	FotoSelfie         *multipart.FileHeader `form:"foto_selfie" binding:"omitempty"`
}

type CreateConsumerInput struct {
	Nik                string
	FullName           string
	LegalName          string
	Email              string
	Password           string
	TempatLahir        string
	TanggalLahir       string
	Gaji               float64
	OverallCreditLimit float64
	FotoKtpPath        string
	FotoSelfiePath     string
}

type UpdateConsumerInput struct {
	FullName           *string  `json:"full_name" validate:"omitempty,min=2"`
	LegalName          *string  `json:"legal_name" validate:"omitempty,min=2"`
	TempatLahir        *string  `json:"tempat_lahir" validate:"omitempty,min=1"`
	TanggalLahir       *string  `json:"tanggal_lahir" validate:"omitempty,datetime=2006-01-02"`
	Gaji               *float64 `json:"gaji" validate:"omitempty,gt=0"`
	OverallCreditLimit *float64 `json:"overall_credit_limit" validate:"omitempty,gte=0"`
	FotoKtp            *string  `json:"foto_ktp" validate:"omitempty,url|uri"`
	FotoSelfie         *string  `json:"foto_selfie" validate:"omitempty,url|uri"`
}
