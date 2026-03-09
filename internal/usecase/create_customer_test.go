package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/internal/usecase"
	"github.com/belayhun-arage/billing-service/test/mocks"
)

func TestCreateCustomer_HappyPath(t *testing.T) {
	repo := &mocks.MockCustomerRepository{
		ExistsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
		CreateFn:        func(_ context.Context, _ *domain.Customer) error { return nil },
	}

	uc := usecase.NewCreateCustomerUsecase(repo)
	c, err := uc.Execute(context.Background(), "Alice", "alice@example.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name != "Alice" {
		t.Errorf("Name = %q, want %q", c.Name, "Alice")
	}
	if c.Email != "alice@example.com" {
		t.Errorf("Email = %q, want %q", c.Email, "alice@example.com")
	}
	if c.ID == "" {
		t.Error("ID must be set")
	}
}

func TestCreateCustomer_DuplicateEmail(t *testing.T) {
	repo := &mocks.MockCustomerRepository{
		ExistsByEmailFn: func(_ context.Context, _ string) (bool, error) { return true, nil },
	}

	uc := usecase.NewCreateCustomerUsecase(repo)
	_, err := uc.Execute(context.Background(), "Alice", "alice@example.com")

	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestCreateCustomer_InvalidDomainInput(t *testing.T) {
	// ExistsByEmail should never be called — domain validation fails first.
	called := false
	repo := &mocks.MockCustomerRepository{
		ExistsByEmailFn: func(_ context.Context, _ string) (bool, error) {
			called = true
			return false, nil
		},
	}

	uc := usecase.NewCreateCustomerUsecase(repo)

	tests := []struct{ name, email string }{
		{"", "alice@example.com"},
		{"Alice", "not-an-email"},
	}
	for _, tc := range tests {
		_, err := uc.Execute(context.Background(), tc.name, tc.email)
		if err == nil {
			t.Errorf("Execute(%q, %q): expected validation error, got nil", tc.name, tc.email)
		}
	}
	if called {
		t.Error("ExistsByEmail must not be called when domain validation fails")
	}
}

func TestCreateCustomer_ExistsByEmailRepoError(t *testing.T) {
	repoErr := errors.New("db connection lost")
	repo := &mocks.MockCustomerRepository{
		ExistsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, repoErr },
	}

	uc := usecase.NewCreateCustomerUsecase(repo)
	_, err := uc.Execute(context.Background(), "Alice", "alice@example.com")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got %v", err)
	}
}

func TestCreateCustomer_CreateRepoError(t *testing.T) {
	repoErr := errors.New("insert failed")
	repo := &mocks.MockCustomerRepository{
		ExistsByEmailFn: func(_ context.Context, _ string) (bool, error) { return false, nil },
		CreateFn:        func(_ context.Context, _ *domain.Customer) error { return repoErr },
	}

	uc := usecase.NewCreateCustomerUsecase(repo)
	_, err := uc.Execute(context.Background(), "Alice", "alice@example.com")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got %v", err)
	}
}
