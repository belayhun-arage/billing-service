package postgres

import (
	"context"

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
	INSERT INTO payments (id, invoice_id, amount, status)
	VALUES ($1,$2,$3,$4)
	`

	args := []any{payment.ID, payment.InvoiceID, payment.Amount, payment.Status}

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = r.db.Exec(ctx, query, args...)
	}

	return err
}
