package usecase

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type CreateMerchantUsecase struct {
	repo domain.MerchantRepository
}

func NewCreateMerchantUsecase(repo domain.MerchantRepository) *CreateMerchantUsecase {
	return &CreateMerchantUsecase{repo: repo}
}

func (u *CreateMerchantUsecase) Execute(ctx context.Context, name, email string) (*domain.Merchant, error) {
	merchant, err := domain.NewMerchant(name, email)
	if err != nil {
		return nil, err
	}
	if err := u.repo.Create(ctx, merchant); err != nil {
		return nil, err
	}
	return merchant, nil
}
