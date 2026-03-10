package mocks

import (
	"context"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// MockCustomerRepository satisfies both domain.CustomerRepository and the
// extended interface required by CreateCustomerUsecase.
type MockCustomerRepository struct {
	CreateFn                  func(ctx context.Context, c *domain.Customer) error
	GetByIDFn                 func(ctx context.Context, id string) (*domain.Customer, error)
	ExistsByEmailForMerchantFn func(ctx context.Context, merchantID, email string) (bool, error)
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

func (m *MockCustomerRepository) ExistsByEmailForMerchant(ctx context.Context, merchantID, email string) (bool, error) {
	if m.ExistsByEmailForMerchantFn != nil {
		return m.ExistsByEmailForMerchantFn(ctx, merchantID, email)
	}
	return false, nil
}
