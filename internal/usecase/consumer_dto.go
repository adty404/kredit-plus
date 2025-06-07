package usecase

type CreateConsumerInput struct {
	Nik            string
	FullName       string
	LegalName      string
	TempatLahir    string
	TanggalLahir   string
	Gaji           float64
	FotoKtpPath    string // Path file KTP yang sudah disimpan
	FotoSelfiePath string // Path file Selfie yang sudah disimpan
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
