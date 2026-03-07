package domain

import (
	"context"
	"time"

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

type InvoiceRepository interface {
	Create(ctx context.Context, tx pgx.Tx, invoice *Invoice) error
	GetByID(ctx context.Context, id string) (*Invoice, error)
}
