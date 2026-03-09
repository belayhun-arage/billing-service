package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

// MockInvoiceRepository is a test double for domain.InvoiceRepository.
type MockInvoiceRepository struct {
	CreateFn  func(ctx context.Context, tx pgx.Tx, invoice *domain.Invoice) error
	GetByIDFn func(ctx context.Context, id string) (*domain.Invoice, error)
}

func (m *MockInvoiceRepository) Create(ctx context.Context, tx pgx.Tx, invoice *domain.Invoice) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, tx, invoice)
	}
	return nil
}

func (m *MockInvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
