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

// Execute revokes the given key. callerCustomerID must match the key's owner;
// returning ErrAPIKeyNotFound in both the not-found and the unauthorized cases
// avoids leaking key existence to unauthorized callers.
func (u *RevokeAPIKeyUsecase) Execute(ctx context.Context, key, callerCustomerID string) error {
	if key == "" {
		return domain.ErrAPIKeyNotFound
	}

	existing, err := u.repo.GetByKey(ctx, key)
	if err != nil {
		return domain.ErrAPIKeyNotFound
	}

	if existing.CustomerID != callerCustomerID {
		return domain.ErrAPIKeyNotFound
	}

	return u.repo.Revoke(ctx, key)
}
