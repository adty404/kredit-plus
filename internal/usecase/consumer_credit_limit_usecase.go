package usecase

import (
	"fmt"
	"github.com/adty404/kredit-plus/internal/domain"
)

type ConsumerCreditLimitUsecase interface {
	CreateConsumerCreditLimit(
		consumerID uint,
		input CreateConsumerCreditLimitInput,
	) (*domain.ConsumerCreditLimit, error)
}

type consumerCreditLimitUsecase struct {
	repo         domain.ConsumerCreditLimitRepository
	consumerRepo domain.ConsumerRepository
}

func NewConsumerCreditLimitUsecase(
	repo domain.ConsumerCreditLimitRepository,
	consumerRepo domain.ConsumerRepository,
) ConsumerCreditLimitUsecase {
	return &consumerCreditLimitUsecase{
		repo:         repo,
		consumerRepo: consumerRepo,
	}
}

func (uc *consumerCreditLimitUsecase) CreateConsumerCreditLimit(
	consumerID uint,
	input CreateConsumerCreditLimitInput,
) (*domain.ConsumerCreditLimit, error) {
	// Validasi 1: Pastikan konsumen ada
	consumer, err := uc.consumerRepo.FindByID(consumerID)
	if err != nil {
		return nil, fmt.Errorf("consumer with id %d not found", consumerID)
	}

	// Validasi 2: Pastikan limit untuk tenor ini belum ada
	_, err = uc.repo.FindByConsumerAndTenor(consumerID, input.TenorMonths)
	if err == nil {
		return nil, fmt.Errorf("credit limit for tenor %d months already exists for this consumer", input.TenorMonths)
	}

	// Validasi 3: Tenor Months harus dalam 1, 2, 3, atau 6 bulan
	allowedTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !allowedTenors[input.TenorMonths] {
		return nil, fmt.Errorf("invalid tenor: %d. allowed tenors are 1, 2, 3, 6", input.TenorMonths)
	}

	// Validasi 4: Pastikan limit per tenor tidak melebihi plafon kredit keseluruhan
	if input.CreditLimit > consumer.OverallCreditLimit {
		return nil, fmt.Errorf(
			"credit limit (%.2f) cannot exceed consumer's overall credit limit (%.2f)",
			input.CreditLimit,
			consumer.OverallCreditLimit,
		)
	}

	limit := &domain.ConsumerCreditLimit{
		ConsumerID:  consumerID,
		TenorMonths: input.TenorMonths,
		CreditLimit: input.CreditLimit,
	}

	if err := uc.repo.Save(limit); err != nil {
		return nil, err
	}

	return limit, nil
}
