package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

// APIKey is a merchant-scoped credential used by the merchant's own backend
// systems to authenticate against this billing API.
type APIKey struct {
	ID         string
	MerchantID string
	Key        string
	Secret     string // raw secret — used as HMAC signing key; stored in DB, shown once
	Label      string // human-readable label, e.g. "production", "dashboard"
	CreatedAt  time.Time
	RevokedAt  *time.Time
}

// NewAPIKey generates a new API key + secret pair scoped to the given merchant.
// The Secret is the raw signing secret and must be returned to the caller once
// — it is stored in the DB and never exposed again.
func NewAPIKey(merchantID, label string) (*APIKey, error) {
	if merchantID == "" {
		return nil, errors.New("merchant_id is required to create an api key")
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
		MerchantID: merchantID,
		Key:        "bk_" + hex.EncodeToString(keyBytes),
		Secret:     hex.EncodeToString(secretBytes),
		Label:      label,
		CreatedAt:  time.Now(),
	}, nil
}

var ErrAPIKeyNotFound = errors.New("api key not found or already revoked")

func (k *APIKey) IsActive() bool {
	return k.RevokedAt == nil
}

type APIKeyRepository interface {
	Create(ctx context.Context, key *APIKey) error
	GetByKey(ctx context.Context, key string) (*APIKey, error)
	Revoke(ctx context.Context, key string) error
}
