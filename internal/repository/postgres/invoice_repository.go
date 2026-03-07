package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/belayhun-arage/billing-service/internal/domain"
)

type InvoiceRepository struct {
	db *pgxpool.Pool
}

func NewInvoiceRepository(db *pgxpool.Pool) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) Create(ctx context.Context, tx pgx.Tx, invoice *domain.Invoice) error {

	query := `
	INSERT INTO invoices (id, customer_id, amount, status, created_at)
	VALUES ($1,$2,$3,$4,$5)
	`

	args := []any{invoice.ID, invoice.CustomerID, invoice.Amount, invoice.Status, invoice.CreatedAt}

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = r.db.Exec(ctx, query, args...)
	}

	return err
}

func (r *InvoiceRepository) GetByID(ctx context.Context, id string) (*domain.Invoice, error) {

	query := `
	SELECT id, customer_id, amount, status, created_at
	FROM invoices
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var invoice domain.Invoice

	err := row.Scan(
		&invoice.ID,
		&invoice.CustomerID,
		&invoice.Amount,
		&invoice.Status,
		&invoice.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &invoice, nil
}