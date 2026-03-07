package messaging

import (
    "context"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
    writer *kafka.Writer
}

// NewKafkaPublisher initializes a Kafka publisher
func NewKafkaPublisher(brokers []string, topic string) *KafkaPublisher {
    return &KafkaPublisher{
        writer: &kafka.Writer{
            Addr:         kafka.TCP(brokers...),
            Topic:        topic,
            BatchTimeout: 500 * time.Millisecond, // controls batching
            RequiredAcks: kafka.RequireAll,
        },
    }
}

// Publish sends an event to Kafka
func (k *KafkaPublisher) Publish(eventType string, payload []byte) error {
    msg := kafka.Message{
        Key:   []byte(eventType),
        Value: payload,
        Time:  time.Now(),
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := k.writer.WriteMessages(ctx, msg)
    if err != nil {
        log.Printf("Kafka publish failed: %v", err)
    }

    return err
}

// Close closes the Kafka writer
func (k *KafkaPublisher) Close() error {
    return k.writer.Close()
}