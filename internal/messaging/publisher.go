package messaging

import "log/slog"

type EventPublisher interface {
	Publish(eventType string, payload []byte) error
}

// NoOpPublisher is used when Kafka is not configured.
type NoOpPublisher struct{ log *slog.Logger }

func NewNoOpPublisher(log *slog.Logger) *NoOpPublisher {
	return &NoOpPublisher{log: log}
}

func (n *NoOpPublisher) Publish(eventType string, _ []byte) error {
	n.log.Debug("outbox event dropped — Kafka not configured", "event_type", eventType)
	return nil
}
