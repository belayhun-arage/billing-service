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
	revokeCalled := false
	repo := &mocks.MockAPIKeyRepository{
		RevokeFn: func(_ context.Context, _ string) error {
			revokeCalled = true
			return nil
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_testkey")

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

	err := uc.Execute(context.Background(), "")

	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Errorf("expected ErrAPIKeyNotFound, got %v", err)
	}
}

func TestRevokeAPIKey_KeyNotFound(t *testing.T) {
	repo := &mocks.MockAPIKeyRepository{
		RevokeFn: func(_ context.Context, _ string) error {
			return domain.ErrAPIKeyNotFound
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_missing")

	if !errors.Is(err, domain.ErrAPIKeyNotFound) {
		t.Errorf("expected ErrAPIKeyNotFound, got %v", err)
	}
}

func TestRevokeAPIKey_RepoError(t *testing.T) {
	repoErr := errors.New("db connection lost")
	repo := &mocks.MockAPIKeyRepository{
		RevokeFn: func(_ context.Context, _ string) error {
			return repoErr
		},
	}

	uc := usecase.NewRevokeAPIKeyUsecase(repo)
	err := uc.Execute(context.Background(), "bk_testkey")

	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got %v", err)
	}
}
