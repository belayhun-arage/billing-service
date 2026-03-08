package usecase

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type CreateAPIKeyUsecase struct {
	repo domain.APIKeyRepository
}

func NewCreateAPIKeyUsecase(repo domain.APIKeyRepository) *CreateAPIKeyUsecase {
	return &CreateAPIKeyUsecase{repo: repo}
}

type APIKeyResult struct {
	ID         string `json:"id"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
	CustomerID string `json:"customer_id"`
}

func (u *CreateAPIKeyUsecase) Execute(ctx context.Context, customerID string) (*APIKeyResult, error) {
	apiKey, err := domain.NewAPIKey(customerID)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &APIKeyResult{
		ID:         apiKey.ID,
		Key:        apiKey.Key,
		Secret:     apiKey.Secret,
		CustomerID: apiKey.CustomerID,
	}, nil
}
