package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, tx pgx.Tx, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, invoice_id, customer_id, amount, status, provider_payment_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	args := []any{
		payment.ID,
		payment.InvoiceID,
		payment.CustomerID,
		payment.Amount,
		payment.Status,
		payment.ProviderPaymentID,
		payment.CreatedAt,
		payment.UpdatedAt,
	}

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = r.db.Exec(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("create payment: %w", err)
	}
	return nil
}
