package usecase

type CreateTransactionInput struct {
	TenorMonths     int     `json:"tenor_months" binding:"required,gt=0"`
	Otr             float64 `json:"otr" binding:"required,gt=0"`
	AdminFee        float64 `json:"admin_fee" binding:"gte=0"`
	UangMuka        float64 `json:"uang_muka" binding:"gte=0"`
	NamaAsset       string  `json:"nama_asset" binding:"required"`
	JenisAsset      string  `json:"jenis_asset" binding:"required"`
	SumberTransaksi string  `json:"sumber_transaksi" binding:"required"` // <-- Field baru ditambahkan
}
