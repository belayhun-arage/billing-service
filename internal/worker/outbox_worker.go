package worker

import (
    "context"
    "log"
    "time"

    "github.com/belayhun-arage/billing-service/internal/messaging"
    "github.com/jackc/pgx/v5/pgxpool"
)

type OutboxWorker struct {
    db        *pgxpool.Pool
    publisher messaging.EventPublisher
}

func NewOutboxWorker(db *pgxpool.Pool, publisher messaging.EventPublisher) *OutboxWorker {
    return &OutboxWorker{db: db, publisher: publisher}
}

func (w *OutboxWorker) Start(ctx context.Context) {
    for {
        rows, err := w.db.Query(ctx, `
            SELECT id, event_type, payload
            FROM outbox_events
            WHERE processed = false
            ORDER BY created_at
            LIMIT 50
            FOR UPDATE SKIP LOCKED
        `)
        if err != nil {
            log.Println("Outbox query failed:", err)
            continue
        }

        for rows.Next() {
            var id string
            var eventType string
            var payload []byte

            if err := rows.Scan(&id, &eventType, &payload); err != nil {
                log.Println(err)
                continue
            }

            if err := w.publisher.Publish(eventType, payload); err != nil {
                log.Println("Failed to publish event, will retry:", err)
                continue
            }

            // mark as processed
            _, err = w.db.Exec(ctx, `
                UPDATE outbox_events
                SET processed = true, processed_at = now()
                WHERE id = $1
            `, id)
            if err != nil {
                log.Println("Failed to mark outbox event processed:", err)
            }
        }

        rows.Close()
        // small delay between batches
        select {
        case <-ctx.Done():
            return
        case <-time.After(1 * time.Second):
        }
    }
}