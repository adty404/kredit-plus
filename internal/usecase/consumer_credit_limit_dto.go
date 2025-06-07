package usecase

type CreateConsumerCreditLimitInput struct {
	TenorMonths int     `json:"tenor_months" binding:"required,gt=0"`
	CreditLimit float64 `json:"credit_limit" binding:"required,gte=0"`
}
