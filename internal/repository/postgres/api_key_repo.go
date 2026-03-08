package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type APIKeyRepository struct {
	db *pgxpool.Pool
}

func NewAPIKeyRepository(db *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, k *domain.APIKey) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO api_keys (id, key, secret, customer_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, k.ID, k.Key, k.Secret, k.CustomerID, k.CreatedAt)
	return err
}

func (r *APIKeyRepository) Revoke(ctx context.Context, key string) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE api_keys SET revoked_at = now()
		WHERE key = $1 AND revoked_at IS NULL
	`, key)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrAPIKeyNotFound
	}
	return nil
}

func (r *APIKeyRepository) GetByKey(ctx context.Context, key string) (*domain.APIKey, error) {
	var k domain.APIKey
	err := r.db.QueryRow(ctx, `
		SELECT id, key, secret, customer_id, created_at, revoked_at
		FROM api_keys
		WHERE key = $1
	`, key).Scan(
		&k.ID,
		&k.Key,
		&k.Secret,
		&k.CustomerID,
		&k.CreatedAt,
		&k.RevokedAt,
	)
	if err != nil {
		return nil, err
	}
	return &k, nil
}
