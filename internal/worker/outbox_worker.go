package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/belayhun-arage/billing-service/internal/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxWorker struct {
	db        *pgxpool.Pool
	publisher messaging.EventPublisher
	log       *slog.Logger
}

func NewOutboxWorker(db *pgxpool.Pool, publisher messaging.EventPublisher, log *slog.Logger) *OutboxWorker {
	return &OutboxWorker{db: db, publisher: publisher, log: log}
}

// Start polls for unprocessed outbox events and publishes them. It respects ctx
// for graceful shutdown and backs off exponentially on repeated errors.
func (w *OutboxWorker) Start(ctx context.Context) {
	backoff := time.Second
	const maxBackoff = 30 * time.Second

	for {
		if err := w.processBatch(ctx); err != nil {
			w.log.Error("outbox batch failed, backing off", "error", err, "backoff", backoff)
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// Reset backoff on success.
		backoff = time.Second

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

// processBatch fetches up to 50 unprocessed events inside a single transaction,
// publishes each, and marks them processed atomically.
// FOR UPDATE SKIP LOCKED prevents concurrent workers from double-processing.
func (w *OutboxWorker) processBatch(ctx context.Context) error {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	rows, err := tx.Query(ctx, `
		SELECT id, event_type, payload
		FROM outbox_events
		WHERE processed = false
		ORDER BY created_at
		LIMIT 50
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return fmt.Errorf("query outbox events: %w", err)
	}

	type outboxEvent struct {
		id        string
		eventType string
		payload   []byte
	}

	// Collect all rows before closing the cursor; pgx cannot interleave
	// reads and writes on the same connection.
	var events []outboxEvent
	for rows.Next() {
		var e outboxEvent
		if err := rows.Scan(&e.id, &e.eventType, &e.payload); err != nil {
			return fmt.Errorf("scan outbox row: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate outbox rows: %w", err)
	}
	rows.Close()

	for _, e := range events {
		if err := w.publisher.Publish(e.eventType, e.payload); err != nil {
			// Log and skip — event stays unprocessed and will be retried next batch.
			w.log.Error("failed to publish event, will retry next batch",
				"event_id", e.id, "event_type", e.eventType, "error", err)
			continue
		}

		if _, err := tx.Exec(ctx, `
			UPDATE outbox_events
			SET processed = true, processed_at = now()
			WHERE id = $1
		`, e.id); err != nil {
			return fmt.Errorf("mark event %s processed: %w", e.id, err)
		}
	}

	return tx.Commit(ctx)
}
