package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Invoice struct {
	ID         string
	CustomerID string
	Amount     int64
	Currency   string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewInvoice validates inputs and returns an Invoice ready to persist.
func NewInvoice(customerID string, amount int64) (*Invoice, error) {
	if customerID == "" {
		return nil, errors.New("customer_id is required")
	}
	if amount <= 0 {
		return nil, errors.New("invoice amount must be greater than zero")
	}
	now := time.Now()
	return &Invoice{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Amount:     amount,
		Currency:   "usd",
		Status:     "pending",
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

type InvoiceRepository interface {
	Create(ctx context.Context, tx pgx.Tx, invoice *Invoice) error
	GetByID(ctx context.Context, id string) (*Invoice, error)
}
