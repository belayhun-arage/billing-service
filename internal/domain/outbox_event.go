package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type OutboxEvent struct {
	ID        string
	EventType string
	Payload   []byte
	CreatedAt time.Time
	Processed bool
}

type OutboxRepository interface {
	Insert(ctx context.Context, tx pgx.Tx, event *OutboxEvent) error
}