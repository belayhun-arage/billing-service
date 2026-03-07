package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/belayhun-arage/billing-service/internal/domain"
)

type LedgerRepository struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) Create(ctx context.Context, tx pgx.Tx, entry *domain.LedgerEntry) error {

	query := `
	INSERT INTO ledger_entries (id, account_id, type, reference_id, amount, description, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	args := []any{entry.ID, entry.AccountID, entry.Type, entry.ReferenceID, entry.Amount, entry.Description, entry.CreatedAt}

	var err error
	if tx != nil {
		_, err = tx.Exec(ctx, query, args...)
	} else {
		_, err = r.db.Exec(ctx, query, args...)
	}

	return err
}
