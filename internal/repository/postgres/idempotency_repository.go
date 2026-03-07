package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)
type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Get(ctx context.Context, key string) ([]byte, int, error) {

	query := `
	SELECT response, status_code
	FROM idempotency_keys
	WHERE key = $1
	`

	var response []byte
	var status int

	err := r.db.QueryRow(ctx, query, key).Scan(&response, &status)
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

	query := `
	INSERT INTO idempotency_keys (key, response, status_code)
	VALUES ($1,$2,$3)
	`

	_, err := r.db.Exec(ctx, query, key, response, status)

	return err
}