package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Merchant is the top-level tenant in the billing service.
// Every API key, customer, invoice, and payment is scoped to a merchant.
type Merchant struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewMerchant(name, email string) (*Merchant, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("merchant name is required")
	}
	if !isValidEmail(email) {
		return nil, errors.New("invalid merchant email address")
	}
	return &Merchant{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		CreatedAt: time.Now(),
	}, nil
}

var ErrMerchantNotFound = errors.New("merchant not found")

type MerchantRepository interface {
	Create(ctx context.Context, merchant *Merchant) error
	GetByID(ctx context.Context, id string) (*Merchant, error)
}
