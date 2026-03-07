package messaging

type EventPublisher interface {
	Publish(eventType string, payload []byte) error
}
