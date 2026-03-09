package mocks

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// MockCustomerRepository satisfies both domain.CustomerRepository and the
// extended interface required by CreateCustomerUsecase (adds ExistsByEmail).
type MockCustomerRepository struct {
	CreateFn        func(ctx context.Context, c *domain.Customer) error
	GetByIDFn       func(ctx context.Context, id string) (*domain.Customer, error)
	ExistsByEmailFn func(ctx context.Context, email string) (bool, error)
}

func (m *MockCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, c)
	}
	return nil
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *MockCustomerRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFn != nil {
		return m.ExistsByEmailFn(ctx, email)
	}
	return false, nil
}
