package domain

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID               string
	MerchantID       string
	Name             string
	Email            string
	StripeCustomerID string // Stripe Customer ID (set after Stripe registration)
	CreatedAt        time.Time
}

// NewCustomer validates inputs and returns a Customer ready to persist.
func NewCustomer(merchantID, name, email string) (*Customer, error) {
	if merchantID == "" {
		return nil, errors.New("merchant_id is required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("customer name is required")
	}
	if !isValidEmail(email) {
		return nil, errors.New("invalid email address")
	}
	return &Customer{
		ID:         uuid.New().String(),
		MerchantID: merchantID,
		Name:       strings.TrimSpace(name),
		Email:      strings.ToLower(strings.TrimSpace(email)),
		CreatedAt:  time.Now(),
	}, nil
}

// isValidEmail checks that the email has exactly one @ with content on both sides
// and a dot in the domain part. No external dependency needed.
func isValidEmail(email string) bool {
	parts := strings.Split(strings.TrimSpace(email), "@")
	return len(parts) == 2 && len(parts[0]) > 0 && strings.Contains(parts[1], ".")
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, id string) (*Customer, error)
}
