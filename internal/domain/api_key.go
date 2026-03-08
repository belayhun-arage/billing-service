package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

type APIKey struct {
	ID         string
	Key        string
	Secret     string // raw secret — used as HMAC signing key
	CustomerID string
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

// NewAPIKey generates a new API key + secret pair for the given customer.
// The Secret is the raw signing secret and must be returned to the caller once
// — it is stored in the DB and never exposed again.
func NewAPIKey(customerID string) (*APIKey, error) {
	if customerID == "" {
		return nil, errors.New("customer_id is required")
	}

	keyBytes := make([]byte, 16)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return nil, err
	}

	return &APIKey{
		ID:         uuid.New().String(),
		Key:        "bk_" + hex.EncodeToString(keyBytes), // e.g. bk_a3f1...
		Secret:     hex.EncodeToString(secretBytes),       // 64-char hex string
		CustomerID: customerID,
		CreatedAt:  time.Now(),
	}, nil
}

func (k *APIKey) IsActive() bool {
	return k.RevokedAt == nil
}

type APIKeyRepository interface {
	Create(ctx context.Context, key *APIKey) error
	GetByKey(ctx context.Context, key string) (*APIKey, error)
}
