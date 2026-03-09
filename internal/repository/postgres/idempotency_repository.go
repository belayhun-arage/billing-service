package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

// hashKey returns a SHA-256 hex digest of the raw idempotency key so that
// user-controlled values are never stored or compared as plaintext.
func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func (r *IdempotencyRepository) Get(ctx context.Context, key string) ([]byte, int, error) {
	var response []byte
	var status int

	err := r.db.QueryRow(ctx, `
		SELECT response, status_code
		FROM idempotency_keys
		WHERE key = $1
	`, hashKey(key)).Scan(&response, &status)
	if err != nil {
		return nil, 0, err
	}

	return response, status, nil
}

func (r *IdempotencyRepository) Save(
	ctx context.Context,
	key string,
	response []byte,
	status int,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO idempotency_keys (key, response, status_code)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO NOTHING
	`, hashKey(key), response, status)
	return err
}
