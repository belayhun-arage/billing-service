package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/belayhun-arage/billing-service/internal/domain"
	"github.com/belayhun-arage/billing-service/internal/usecase"
	"github.com/belayhun-arage/billing-service/test/mocks"
)

func TestRevokeAPIKey_HappyPath(t *testing.T) {
	const (
		keyStr     = "bk_testkey"
		customerID = "cust-123"
	)

	activeKey := &domain.APIKey{
		ID:         "key-id-1",
		Key:        keyStr,
		Secret:     "somesecret",
		CustomerID: customerID,
	}

	revokeCalled := false
	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, k string) (*domain.APIKey, error) {
			if k == keyStr {
				return activeKey, nil
			}
			return nil, domain.ErrAPIKeyNotFound
		},
		RevokeFn: func(_ context.Context, _ string) error {
			revokeCalled = true
			return nil
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), keyStr, customerID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !revokeCalled {
		t.Error("expected repo.Revoke to be called")
	}
}

func TestRevokeAPIKey_EmptyKey(t *testing.T) {
	repo := &mocks.MockAPIKeyRepository{}
	uc := usecase.NewRevokeAPIKeyUsecase(repo)

	err := uc.Execute(context.Background(), "", "cust-123")

	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Errorf("expected ErrAPIKeyNotFound, got %v", err)
	}
}

func TestRevokeAPIKey_KeyNotFound(t *testing.T) {
	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, _ string) (*domain.APIKey, error) {
			return nil, domain.ErrAPIKeyNotFound
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_missing", "cust-123")

	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Errorf("expected ErrAPIKeyNotFound, got %v", err)
	}
}

func TestRevokeAPIKey_WrongOwner(t *testing.T) {
	// Key belongs to cust-A; caller claims to be cust-B.
	key := &domain.APIKey{
		ID:         "key-id-1",
		Key:        "bk_testkey",
		Secret:     "somesecret",
		CustomerID: "cust-A",
	}

	revokeCalled := false
	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, _ string) (*domain.APIKey, error) {
			return key, nil
		},
		RevokeFn: func(_ context.Context, _ string) error {
			revokeCalled = true
			return nil
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_testkey", "cust-B")

	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Errorf("expected ErrAPIKeyNotFound for unauthorized caller, got %v", err)
	}
	if revokeCalled {
		t.Error("repo.Revoke must not be called when the caller doesn't own the key")
	}
}

func TestRevokeAPIKey_RepoRevokeError(t *testing.T) {
	repoErr := errors.New("db error")
	key := &domain.APIKey{
		ID:         "key-id-1",
		Key:        "bk_testkey",
		CustomerID: "cust-123",
	}

	repo := &mocks.MockAPIKeyRepository{
		GetByKeyFn: func(_ context.Context, _ string) (*domain.APIKey, error) {
			return key, nil
		},
		RevokeFn: func(_ context.Context, _ string) error {
			return repoErr
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_testkey", "cust-123")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got %v", err)
	}
}
