package usecase

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type RevokeAPIKeyUsecase struct {
	repo domain.APIKeyRepository
}

func NewRevokeAPIKeyUsecase(repo domain.APIKeyRepository) *RevokeAPIKeyUsecase {
	return &RevokeAPIKeyUsecase{repo: repo}
}

// Execute revokes the given service-level API key.
// Returns ErrAPIKeyNotFound if the key does not exist or is already revoked.
func (u *RevokeAPIKeyUsecase) Execute(ctx context.Context, key string) error {
	if key == "" {
		return domain.ErrAPIKeyNotFound
	}
	return u.repo.Revoke(ctx, key)
}
