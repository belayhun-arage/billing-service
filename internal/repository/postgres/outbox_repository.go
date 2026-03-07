package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/belayhun-arage/billing-service/internal/domain"
)

type OutboxRepository struct {
	db *pgxpool.Pool
}

func (r *OutboxRepository) Insert(
	ctx context.Context,
	tx pgx.Tx,
	event *domain.OutboxEvent,
) error {

	query := `
	INSERT INTO outbox_events (id, event_type, payload, created_at)
	VALUES ($1,$2,$3,$4)
	`

	_, err := tx.Exec(
		ctx,
		query,
		event.ID,
		event.EventType,
		event.Payload,
		event.CreatedAt,
	)

	return err
}