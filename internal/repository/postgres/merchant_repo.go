package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type MerchantRepository struct {
	db *pgxpool.Pool
}

func NewMerchantRepository(db *pgxpool.Pool) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Create(ctx context.Context, m *domain.Merchant) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO merchants (id, name, email, created_at)
		VALUES ($1, $2, $3, $4)
	`, m.ID, m.Name, m.Email, m.CreatedAt)
	if err != nil {
		return fmt.Errorf("create merchant: %w", err)
	}
	return nil
}

func (r *MerchantRepository) GetByID(ctx context.Context, id string) (*domain.Merchant, error) {
	var m domain.Merchant
	err := r.db.QueryRow(ctx, `
		SELECT id, name, email, created_at
		FROM merchants
		WHERE id = $1
	`, id).Scan(&m.ID, &m.Name, &m.Email, &m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get merchant %s: %w", id, err)
	}
	return &m, nil
}
