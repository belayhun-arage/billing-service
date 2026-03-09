package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/internal/usecase"
	"github.com/belayhun-arage/billing-service/test/mocks"
)

func TestCreateInvoice_HappyPath(t *testing.T) {
	repo := &mocks.MockInvoiceRepository{
		CreateFn: func(_ context.Context, _ pgx.Tx, inv *domain.Invoice) error { return nil },
	}

	uc := usecase.NewCreateInvoiceUsecase(repo)
	inv, err := uc.Execute(context.Background(), "cust-123", 9999)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.CustomerID != "cust-123" {
		t.Errorf("CustomerID = %q, want %q", inv.CustomerID, "cust-123")
	}
	if inv.Amount != 9999 {
		t.Errorf("Amount = %d, want %d", inv.Amount, 9999)
	}
	if inv.ID == "" {
		t.Error("ID must be set")
	}
}

func TestCreateInvoice_InvalidDomainInput(t *testing.T) {
	called := false
	repo := &mocks.MockInvoiceRepository{
		CreateFn: func(_ context.Context, _ pgx.Tx, _ *domain.Invoice) error {
			called = true
			return nil
		},
	}

	uc := usecase.NewCreateInvoiceUsecase(repo)

	tests := []struct {
		desc       string
		customerID string
		amount     int64
	}{
		{"empty customer_id", "", 5000},
		{"zero amount", "cust-123", 0},
		{"negative amount", "cust-123", -1},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tc.customerID, tc.amount)
			if err == nil {
				t.Errorf("expected validation error, got nil")
			}
		})
	}

	if called {
		t.Error("repo.Create must not be called when domain validation fails")
	}
}

func TestCreateInvoice_RepoError(t *testing.T) {
	repoErr := errors.New("insert failed")
	repo := &mocks.MockInvoiceRepository{
		CreateFn: func(_ context.Context, _ pgx.Tx, _ *domain.Invoice) error { return repoErr },
	}

	uc := usecase.NewCreateInvoiceUsecase(repo)
	_, err := uc.Execute(context.Background(), "cust-123", 5000)

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got %v", err)
	}
}
