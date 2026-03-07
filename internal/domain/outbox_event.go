package domain

import "time"

type OutboxEvent struct {
	ID        string
	EventType string
	Payload   []byte
	CreatedAt time.Time
	Processed bool
}