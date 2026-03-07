package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Payment struct {
	ID         string
	InvoiceID  string
	CustomerID string
	Amount     int64
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
