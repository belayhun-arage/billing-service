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

// APIKeyResult carries the newly created key back to the caller.
// Secret is only available here — it is not exposed again after creation.
type APIKeyResult struct {
	ID     string
	Key    string
	Secret string
	Label  string
}

// Execute creates a new service-level API key with the given label.
// label is optional and used only for human identification (e.g. "production").
func (u *CreateAPIKeyUsecase) Execute(ctx context.Context, label string) (*APIKeyResult, error) {
	apiKey, err := domain.NewAPIKey(label)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &APIKeyResult{
		ID:     apiKey.ID,
		Key:    apiKey.Key,
		Secret: apiKey.Secret,
		Label:  apiKey.Label,
	}, nil
}
