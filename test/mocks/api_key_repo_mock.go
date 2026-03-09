package mocks

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// MockAPIKeyRepository is a test double for domain.APIKeyRepository.
// Set each Fn field to control behaviour per test case.
type MockAPIKeyRepository struct {
	GetByKeyFn func(ctx context.Context, key string) (*domain.APIKey, error)
	CreateFn   func(ctx context.Context, key *domain.APIKey) error
	RevokeFn   func(ctx context.Context, key string) error
}

func (m *MockAPIKeyRepository) GetByKey(ctx context.Context, key string) (*domain.APIKey, error) {
	if m.GetByKeyFn != nil {
		return m.GetByKeyFn(ctx, key)
	}
	return nil, domain.ErrAPIKeyNotFound
}

func (m *MockAPIKeyRepository) Create(ctx context.Context, key *domain.APIKey) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, key)
	}
	return nil
}

func (m *MockAPIKeyRepository) Revoke(ctx context.Context, key string) error {
	if m.RevokeFn != nil {
		return m.RevokeFn(ctx, key)
	}
	return nil
}
