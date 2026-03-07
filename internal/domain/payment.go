package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Payment struct {
	ID                string
	InvoiceID         string
	CustomerID        string
	Amount            int64
	Status            string
	ProviderPaymentID string // Stripe PaymentIntent ID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ChargeResult is returned by PaymentProcessor.Charge.
type ChargeResult struct {
	ProviderPaymentID string
	Status            string
}

// PaymentProcessor is the interface for external payment providers (e.g. Stripe).
// BillingService depends on this interface, not on a concrete SDK type.
type PaymentProcessor interface {
	Charge(providerCustomerID string, amount int64, currency string) (*ChargeResult, error)
	Refund(providerPaymentID string) error
}

type LedgerEntry struct {
	ID          string
	AccountID   string
	Type        string
	ReferenceID string
	Amount      int64
	Description string
	CreatedAt   time.Time
}

type LedgerRepository interface {
	Create(ctx context.Context, tx pgx.Tx, entry *LedgerEntry) error
}

type PaymentRepository interface {
	Create(ctx context.Context, tx pgx.Tx, payment *Payment) error
}
